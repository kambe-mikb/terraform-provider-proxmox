package proxmox

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	pxapi "github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/Telmate/terraform-provider-proxmox/v2/proxmox/Internal/pxapi/guest/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// using a global variable here so that we have an internally accessible
// way to look into our own resource definition. Useful for dynamically doing typecasts
// so that we can print (debug) our ResourceData constructs
var sdnZoneResourceDef *schema.Resource

func resourceSdnZone() *schema.Resource {
	SdnZoneResourceDef = &schema.Resource{
		Create:        resourceSdnZoneCreate,
		Read:          resourceSdnZoneRead,
		Update:        resourceSdnZoneUpdate,
		DeleteContext: resourceSdnZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type: schema.TypeString,
				Required: true,
			},
			"zone": {
				Type: schema.TypeString,
				Required: true,
			},
			"advertise_subnets": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"bridge": {
				Type: schema.TypeString,
				Optional: true,
			},
			"bridge_disable_mac_learning": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"controller": {
				Type: schema.TypeString,
				Optional: true,
			},
			"disable_arp_nd_suppression": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"dns": {
				Type: schema.TypeString,
				Optional: true,
			},
			"dnsZone": {
				Type: schema.TypeString,
				Optional: true,
			},
			"dp_id":{
				Type: schema.TypeInt,
				Optional: true,
			},
			"exitnodes": {
				Type: schema.TypeString,
				Optional: true,
			},
			"exitnodes_local_routing": {
				Type: schema.TypeString,
				Optional: true,
			},
			"exitnodes_primary": {
				Type: schema.TypeString,
				Optional: true,
			},
			"ipam": {
				Type: schema.TypeString,
				Optional: true,
			},
			"mac": {
				Type: schema.TypeString,
				Optional: true,
			},
			"mtu": {
				Type: schema.TypeInt,
				Optional: true,
			},
			"nodes": {
				Type: schema.TypeString,
				Optional: true,
			},
			"peers": {
				Type: schema.TypeString,
				Optional: true,
			},
			"reversedns": {
				Type: schema.TypeString,
				Optional: true,
			},
			"rt_import": {
				Type: schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type: schema.TypeInt,
				Optional: true,
			},
			"vlan_protocol": {
				Type: schema.TypeString,
				Optional: true,
			},
			"vrf_vxlan": {
				Type: schema.TypeString,
				Optional: true,
			},
			"delete": {
				Type: schema.TypeString,
				Optional: true,
			},
			"digest": {
				Type: schema.TypeString,
				Optional: true,
			},
		},
		Timeouts: resourceTimeouts(),
	}

	return SdnZoneResourceDef
}
		
func resourceSdnZoneCreate(d *schema.ResourceData, meta interface{}) error {
	pconf := meta.(*providerConfiguration)

	lock := pmParallelBegin(pconf)
	defer lock.unlock()

	client := pconf.Client

	config := pxapi.ConfigSDNZone{}
	config.Type = d.Get("type").(string)
	config.Zone = d.Get("zone").(string)
	config.AdvertiseSubnets = d.Get("advertise_subnets").(bool)
	config.Bridge = d.Get("bridge").(string)
	config.BridgeDisableMacLearning = d.Get("bridge_disable_mac_learning").(bool)
	config.Controller = d.Get("controller").(string)
	config.DisableARPNDSuppression = d.Get("disable_arp_nd_suppression").(bool)
	config.DNS = d.Get("dns").(string)
	config.DNSZone = d.Get("dnszone").(string)
	config.DPID = d.Get("dp_id").(int)
	config.ExitNodes = d.Get("exitnodes").(string)
	config.ExitNodesLocalRouting = d.Get("exitnodes_local_routing").(bool)
	config.ExitNodesPrimary = d.Get("exitnodes_primary").(string)
	config.IPAM = d.Get("ipam").(string)
	config.MAC = d.Get("mac").(string)
	config.MTU = d.Get("mtu").(string)
	config.Nodes = d.Get("nodes").(string)
	config.Peers = d.Get("peers").(string)
	config.ReverseDNS = d.Get("reversedns").(string)
	config.RTImport = d.Get("rt_import").(string)
	config.Tag = d.Get("tag").(int)
	config.VlanProtocol = d.Get("vlan_protocol").(string)
	config.VrfVxlan = d.Get("vrf_vxlan").(string)
	
	zid = d.Get("zone").(string)
	d.SetId(zid)
	
	err = config.CreateWithValidate(zid, client)
	if err != nil {
		return err
	}
	
	lock.unlock()
	return resourceSdnZoneRead(d, meta)
}

func resourceSdnZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	pconf := meta.(*providerConfiguration)
	lock := pmParallelBegin(pconf)
	defer lock.unlock()

	client := pconf.Client
	zid = d.Get("zone").(string)

	zoneExists, err := pxapi.CheckSDNZoneExistance(zid)
	if err != nil {
		return
	}
	if !zoneExists {
		return fmt.Error("zone %s could not be found", zid)
	}

	c, err := client.getSDNZone(zid)
	if err != nil {
		return
	}
	config, err = pxapi.NewConfigSDNZoneFromJson(c)
	if err != nil {
		return err
	}

	config.Type = d.Get("type").(string)
	config.Zone = d.Get("zone").(string)
	config.AdvertiseSubnets = d.Get("advertise_subnets").(bool)
	config.Bridge = d.Get("bridge").(string)
	config.BridgeDisableMacLearning = d.Get("bridge_disable_mac_learning").(bool)
	config.Controller = d.Get("controller").(string)
	config.DisableARPNDSuppression = d.Get("disable_arp_nd_suppression").(bool)
	config.DNS = d.Get("dns").(string)
	config.DNSZone = d.Get("dnszone").(string)
	config.DPID = d.Get("dp_id").(int)
	config.ExitNodes = d.Get("exitnodes").(string)
	config.ExitNodesLocalRouting = d.Get("exitnodes_local_routing").(bool)
	config.ExitNodesPrimary = d.Get("exitnodes_primary").(string)
	config.IPAM = d.Get("ipam").(string)
	config.MAC = d.Get("mac").(string)
	config.MTU = d.Get("mtu").(string)
	config.Nodes = d.Get("nodes").(string)
	config.Peers = d.Get("peers").(string)
	config.ReverseDNS = d.Get("reversedns").(string)
	config.RTImport = d.Get("rt_import").(string)
	config.Tag = d.Get("tag").(int)
	config.VlanProtocol = d.Get("vlan_protocol").(string)
	config.VrfVxlan = d.Get("vrf_vxlan").(string)

	err = config.UpdateWithValidate(zid, client)
	if err != nil {
		return err
	}

	lock.unlock()
	return resourceSdnZoneRead(d, meta)
}

func resourceSdnZoneRead(d *schema.ResourceData, meta interface{}) error {
	pconf := meta.(*providerConfiguration)
	lock := pmParallelBegin(pconf)
	defer lock.unlock()
	client := pconf.Client
	zid = d.Get("zone").(string)

	zoneExists, err := pxapi.CheckSDNZoneExistance(zid)
	if err != nil {
		return
	}
	if !zoneExists {
		return fmt.Error("zone %s could not be found", zid)
	}

	c, err := client.getSDNZone(zid)
	if err != nil {
		return
	}
	config, err = pxapi.NewConfigSDNZoneFromJson(c)
	if err != nil {
		return err
	}

	d.Set("type") = config.Type
	d.Set("zone") = config.Zone
	d.Set("advertise_subnets") = config.AdvertiseSubnets
	d.Set("bridge") = config.Bridge
	d.Set("bridge_disable_mac_learning") = config.BridgeDisableMacLearning
	d.Set("controller") = config.Controller
	d.Set("disable_arp_nd_suppression") = config.DisableARPNDSuppression
	d.Set("dns") = config.DNS
	d.Set("dnszone") = config.DNSZone
	d.Set("dp_id") = config.DPID
	d.Set("exitnodes") = config.ExitNodes
	d.Set("exitnodes_local_routing") = config.ExitNodesLocalRouting
	d.Set("exitnodes_primary") = config.ExitNodesPrimary
	d.Set("ipam")= config.IPAM
	d.Set("mac") = config.MAC
	d.Set("mtu") = config.MTU
	d.Set("nodes") = config.Nodes
	d.Set("peers") = config.Peers
	d.Set("reversedns") = config.ReverseDNS
	d.Set("rt_import") = config.RTImport
	d.Set("tag") = config.Tag
	d.Set("vlan_protocol") = config.VlanProtocol
	d.Set("vrf_vxlan") = config.VrfVxlan
	
	return nil
}

func resourceSdnZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pconf := meta.(*providerConfiguration)
	lock := pmParallelBegin(pconf)
	defer lock.unlock()

	client := pconf.Client
	zid = d.Get("zone").(string)

	_, err = client.DeleteSDNZone(zid)
	return diag.FromErr(err)
}