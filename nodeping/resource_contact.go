package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func resourceContact() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContactCreate,
		ReadContext:   resourceContactRead,
		UpdateContext: resourceContactUpdate,
		DeleteContext: resourceContactDelete,
		Schema: map[string]*schema.Schema{
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true},
			"name":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"custrole":    &schema.Schema{Type: schema.TypeString, Optional: true},
			"addresses": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":            &schema.Schema{Type: schema.TypeString, Computed: true},
						"address":       &schema.Schema{Type: schema.TypeString, Required: true},
						"type":          &schema.Schema{Type: schema.TypeString, Required: true},
						"suppressup":    &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressdown":  &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressfirst": &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressall":   &schema.Schema{Type: schema.TypeBool, Optional: true},
						// webhooks related attributes
						"action":       &schema.Schema{Type: schema.TypeString, Optional: true},
						"data":         &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString},
						"headers":      &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString},
						"querystrings": &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString},
						// pushover attributes
						"priority": &schema.Schema{Type: schema.TypeInt, Optional: true},
					},
				},
			},
		},
	}
}

func getContactFromSchema(d *schema.ResourceData) *nodeping_api_client.Contact {
	var contact nodeping_api_client.Contact
	contact.ID = d.Id()
	contact.CustomerId = d.Get("customer_id").(string)
	contact.Name = d.Get("name").(string)
	contact.Custrole = d.Get("custrole").(string)

	addrs := d.Get("addresses").(*schema.Set).List()
	addresses := make(map[string]nodeping_api_client.Address)
	newAddresses := make([]nodeping_api_client.Address, 0)
	for _, addr := range addrs {
		a := addr.(map[string]interface{})

		// get address Id (if present)
		addressId := a["id"].(string)

		// convert "data", "headers" and "querystrings" from interface{}
		// to map[string]string
		data := make(map[string]string)
		for key, val := range a["data"].(map[string]interface{}) {
			data[key] = val.(string)
		}
		headers := make(map[string]string)
		for key, val := range a["headers"].(map[string]interface{}) {
			headers[key] = val.(string)
		}
		querystrings := make(map[string]string)
		for key, val := range a["querystrings"].(map[string]interface{}) {
			querystrings[key] = val.(string)
		}

		address := nodeping_api_client.Address{
			ID:            a["id"].(string),
			Address:       a["address"].(string),
			Type:          a["type"].(string),
			Suppressup:    a["suppressup"].(bool),
			Suppressdown:  a["suppressdown"].(bool),
			Suppressfirst: a["suppressfirst"].(bool),
			Suppressall:   a["suppressall"].(bool),
			Action:        a["action"].(string),
			Data:          data,
			Headers:       headers,
			Querystrings:  querystrings,
			Priority:      a["priority"].(int),
		}

		// addresses that have an id go to addresses, the ones that don't,
		// go to new addresses
		if len(addressId) > 0 {
			addresses[addressId] = address
		} else {
			newAddresses = append(newAddresses, address)
		}
	}
	contact.Addresses = addresses
	contact.NewAddresses = newAddresses

	return &contact
}

func resourceContactCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	contact := getContactFromSchema(d)

	savedContact, err := client.CreateContact(contact)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(savedContact.ID)
	return resourceContactRead(ctx, d, m)
}

func resourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	contact, err := client.GetContact(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(contact.ID)
	d.Set("customer_id", contact.CustomerId)
	d.Set("name", contact.Name)
	d.Set("custrole", contact.Custrole)

	addresses := flattenAddresses(&contact.Addresses)
	if err := d.Set("addresses", addresses); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceContactUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	contact := getContactFromSchema(d)

	_, err := client.UpdateContact(contact)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceContactRead(ctx, d, m)
}

func resourceContactDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	err := client.DeleteContact(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	var diags diag.Diagnostics
	return diags
}
