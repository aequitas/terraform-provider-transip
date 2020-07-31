package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"testing"
)

const testSSHKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC/FGVrpv69Tojj67jDr2ny2A2wFNv9EmLiFKanBRbczEgNcEWx+JQ4j5kOzjBNJNFE/9z51me1XnNYvd7IPZxDTY4+a35Q7KVPnSBp6Neroph5vuDJjBa8vA+wZY0kAXkdwAkHSSGc6WTVKSdMl2JZgvw3L/TYJ6Bql3OlOdXTu4qVI+W591/P6XSejv5UbwGEGTwz1LwyGoKFYZgO3wzOjYlgYF8oSODmhRKDns2TVCXMPtQa+AwypL7lC5IRTFKvD2rFJZgQQ+f8firnY9qx5bpMDtkOqGMFJwV0u+NpChr2VPSLN7okXRrPtDGEvIDAqosvSyBfHmGuebk3scTV test@example.com`

func TestAccTransipResourceSSHKey(t *testing.T) {
	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	resource "transip_sshkey" "test" {
		description = "test-%d"
    key         = "%s"
	}
	`, timestamp, testSSHKey)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_sshkey.test", "key", testSSHKey),
					resource.TestCheckResourceAttr("transip_sshkey.test", "md5_fingerprint", "c4:b9:ec:44:ca:ee:81:98:88:59:73:d7:2c:e7:57:59"),
				),
			},
		},
	})
}
