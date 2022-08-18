terraform {
  required_providers {
    grid = {
      source = "threefoldtech/grid"
    }
  }
}

provider "grid" {
}

resource "grid_network" "botnetwork" {
    nodes = [2]
    ip_range = "10.1.0.0/16"
    name = "bot_network"
    add_wg_access = true
}
resource "grid_deployment" "botVm" {
  node = 2
  network_name = grid_network.botnetwork.name
  ip_range = lookup(grid_network.botnetwork.nodes_ip_range, 2, "")
  vms {
    name = "tfgrid_bot"
    flist = "https://hub.grid.tf/helmy.3bot/muhammadhelmy-bot-latest.flist"
    cpu = 2 
    publicip = true
    memory = 4096
    env_vars = {
      SSH_KEY = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDZjOgX91GaRslUCBWUdzavRz36fwwSyqH6Ncfbb6/VfjNeG8ZurnSAWe7rvzEz12CBRu7YkwnjuiShafvTUHzWJ5jFrrpyX/e6Q22pnnUJZz1oNk7HLmqcjKBsJc2xuuF+n6SQPbvJ6U8/NuFAnwFvLF2ewhEdzwHSBQl+bAmt7bfmbB6lpGGEBMrdGyczeE+/uk2fPMzifMjj+th37j1CqxuL/ELdcFS+399ROp0IvwrxOSjzPR5PLrAuR8UxapeXnvwwQ/LvIi5usL/l9WCr8XJhaqrhYgYxfUx12Dzb92cmdKY2FoIvvManVmWkhUuAcbHTVXa3xElIX7ru7xtSeY8g0tGllgDZJRWf9kOqWwe4Ee7ujQ/79b9nLKxeAuMqm8QVJi6c4Olv3dJcZXydrHEpIZAKMrVeJtKl2vfIJNE7AA5muY5sEjQbGEaQyaNB/FDzKhxXIq/J4iaDcGR1C4ytxFskbLDFKJQqxq1L36cE5GyRoFkkycQ9rMTzdM0= helmy@helmy"
    }
    planetary = true
  }
  
}

output "public_ip" {
    value = grid_deployment.botVm.vms[0].computedip
}
