package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/transip/gotransip/v6/ipaddress"
	"github.com/transip/gotransip/v6/vps"
)

// Transfrom the terraform rules to the API rules (FirewallRule)
func vpsFirewallRulesExpand(stateRules []interface{}) ([]vps.FirewallRule, error) {
	rules := make([]vps.FirewallRule, len(stateRules))
	for i, rule := range stateRules {
		r, err := vpsFirewallRuleExpand(rule)
		if err != nil {
			return nil, err
		}
		rules[i] = *r
	}

	return rules, nil
}

// Transfrom the terraform rule to the API rule (FirewallRule)
func vpsFirewallRuleExpand(i interface{}) (*vps.FirewallRule, error) {
	// Top level rule is <string> = <value>
	rawRule := i.(map[string]interface{})
	log.Printf("%+v\n", rawRule)

	// Parse port (can be X or X-Y, if range, split it)
	var startPort, endPort int
	rawPorts := strings.Split(rawRule["port"].(string), "-")

	if len(rawPorts) < 2 {
		startPort, _ = strconv.Atoi(rawPorts[0])
		endPort = startPort
	} else {
		startPort, _ = strconv.Atoi(rawPorts[0])
		endPort, _ = strconv.Atoi(rawPorts[1])
	}

	// Parse whitelist IP addresses
	rawAddresses := rawRule["whitelist"].([]interface{})
	ipAdresses := make([]ipaddress.IPRange, len(rawAddresses))
	for i, ip := range rawAddresses {
		// Parse the IP
		_, ipNetwork, err := net.ParseCIDR(ip.(string))

		// Check if we had some trouble parsing
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q: %s", ip, err)
		}

		ipAdresses[i].IPNet = *ipNetwork
	}

	// Create the rule object
	rule := &vps.FirewallRule{
		Description: rawRule["description"].(string),
		StartPort:   startPort,
		EndPort:     endPort,
		Protocol:    rawRule["protocol"].(string),
		Whitelist:   ipAdresses,
	}

	return rule, nil
}

// Transform the API rule (FirewallRule) to terraform rule
func vpsFirewallRulesFlatten(rules []vps.FirewallRule) []interface{} {

	// Flatten each rule
	stateRules := make([]interface{}, len(rules))
	for i, rule := range rules {
		stateRules[i] = vpsFirewallRuleFlatten(&rule)
	}

	// Return the set
	return stateRules
}

// Transform the API rule (FirewallRule) to terraform rule
func vpsFirewallRuleFlatten(rule *vps.FirewallRule) map[string]interface{} {

	// Parse the port
	var port string
	if rule.StartPort != rule.EndPort {
		port = fmt.Sprintf("%d-%d", rule.StartPort, rule.EndPort)
	} else {
		port = strconv.Itoa(rule.StartPort)
	}

	// Parse the IP addresses
	ipAdresses := make([]interface{}, len(rule.Whitelist))
	for i, ip := range rule.Whitelist {

		ipAdresses[i] = ip.IPNet.String()
	}

	// Build
	res := map[string]interface{}{
		"description": rule.Description,
		"protocol":    rule.Protocol,
		"port":        port,
		"whitelist":   ipAdresses,
	}
	return res
}
