package nodeping_api_client

import (
	"encoding/json"
)

type Address struct {
	ID            string `json:"id,omitempty"`
	Address       string `json:"address"`
	Type          string `json:"type"`
	Suppressup    bool   `json:"suppressup"`
	Suppressdown  bool   `json:"suppressdown"`
	Suppressfirst bool   `json:"suppressfirst"`
	Suppressall   bool   `json:"suppressall"`
	// webhook attribures
	Action       string            `json:"action,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Querystrings map[string]string `json:"querystrings,omitempty"`
	// pushover attributes
	Priority int `json:"priority"`
}

type Check struct {
	ID            string                    `json:"_id,omitempty"`
	Rev           string                    `json:"_rev,omitempty"`
	Label         string                    `json:"label,omitempty"`
	Type          string                    `json:"type,omitempty"`
	CustomerId    string                    `json:"customer_id,omitempty"`
	Description   string                    `json:"description,omitempty"`
	HomeLoc       interface{}               `json:"homeloc"`
	Interval      int                       `json:"interval,omitempty"`
	Status        string                    `json:"status,omitempty"`
	Enable        string                    `json:"enable,omitempty"`
	Public        bool                      `json:"public"`
	Notifications []map[string]Notification `json:"notifications,omitempty"`
	Parameters    map[string]interface{}    `json:"parameters,omitempty"`
	Runlocations  interface{}               `json:"runlocations,omitempty"`
	Created       int                       `json:"created,omitempty"`
	Modified      int                       `json:"modified,omitempty"`
	Queue         interface{}               `json:"queue,omitempty"`
	Uuid          string                    `json:"uuid,omitempty"`
	State         interface{}               `json:"state,omitempty"`
	Firstdown     int                       `json:"firstdown,omitempty"`
}

type CheckUpdate struct { // used for PUT and POST requests.
	/*
		Since checks require a different structure for PUT and POST request,
		compared	to the one received from GET requests, this is a separate struct
		for creating and updating checks.
	*/
	ID            string                    `json:"_id,omitempty"`
	Label         string                    `json:"label,omitempty"`
	CustomerId    string                    `json:"customerid,omitempty"`
	Type          string                    `json:"type,omitempty"`
	Target        string                    `json:"target,omitempty"`
	Interval      int                       `json:"interval,omitempty"`
	Enable        string                    `json:"enabled,omitempty"` // Note this is called `enable` on GET responses
	Public        string                    `json:"public,omitempty"`
	RunLocations  []string                  `json:"runlocations,omitempty"`
	HomeLoc       interface{}               `json:"homeloc"`
	Notifications []map[string]Notification `json:"notifications,omitempty"`
	Threshold     int                       `json:"threshold,omitempty"`
	Sens          int                       `json:"sens,omitempty"`
	Dep           string                    `json:"dep,omitempty"`
	Description   string                    `json:"description,omitempty"`
	// the following are only relevant for certain types
	CheckToken     string                 `json:"checktoken,omitempty"`
	ClientCert     interface{}            `json:"clientcert,omitempty"`
	ContentString  string                 `json:"contentstring,omitempty"`
	Dohdot         string                 `json:"dohdot,omitempty"`
	DnsType        string                 `json:"dnstype,omitempty"`
	DnsToResolve   string                 `json:"dnstoresolve,omitempty"`
	DnsSection     string                 `json:"dnssection,omitempty"`
	Dnsrd          bool                   `json:"dnsrd,omitempty"`
	Transport      string                 `json:"transport,omitempty"`
	Follow         bool                   `json:"follow"`
	Email          string                 `json:"email,omitempty"`
	Port           int                    `json:"port,omitempty"`
	Username       string                 `json:"username,omitempty"`
	Password       string                 `json:"password,omitempty"`
	Secure         string                 `json:"secure,omitempty"`
	Verify         bool                   `json:"verify,omitempty"`
	Ignore         string                 `json:"ignore,omitempty"`
	Invert         bool                   `json:"invert"`
	WarningDays    int                    `json:"warningdays,omitempty"`
	Fields         map[string]CheckField  `json:"fields,omitempty"`
	Postdata       string                 `json:"postdata,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
	ReceiveHeaders map[string]interface{} `json:"receiveheaders,omitempty"`
	SendHeaders    map[string]interface{} `json:"sendheaders,omitempty"`
	Edns           map[string]interface{} `json:"edns,omitempty"`
	Method         string                 `json:"method,omitempty"`
	Statuscode     int                    `json:"statuscode,omitempty"`
	Ipv6           bool                   `json:"ipv6"`
	Regex          bool                   `json:"regex,omitempty"`
	ServerName     string                 `json:"servername,omitempty"`
	Snmpv          string                 `json:"snmpv,omitempty"`
	Snmpcom        string                 `json:"snmpcom,omitempty"`
	VerifyVolume   bool                   `json:"verifyvolume,omitempty"`
	VolumeMin      int                    `json:"volumemin,omitempty"`
	WhoisServer    string                 `json:"whoisserver,omitempty"`
}

type Contact struct {
	/*
		Note that "addresses" can't be omitted from json, even if it's empty, as
		an empty "addresses" map might mean that some addresses should be
		removed.
	*/
	ID           string             `json:"_id,omitempty"`
	Type         string             `json:"type,omitempty"`
	CustomerId   string             `json:"customer_id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Custrole     string             `json:"custrole,omitempty"`
	Addresses    map[string]Address `json:"addresses"`
	NewAddresses []Address          `json:"newaddresses,omitempty"`
}

func (c *Contact) MarshalJSONForCreate() ([]byte, error) {
	/*
		When calling API to create a new contract, passed json object is not
		allowed to have "addresses" field, and doesn't need the "id" field.
	*/
	return json.Marshal(struct {
		CustomerId   string    `json:"customerid,omitempty"`
		Name         string    `json:"name,omitempty"`
		Custrole     string    `json:"custrole,omitempty"`
		NewAddresses []Address `json:"newaddresses,omitempty"`
	}{c.CustomerId, c.Name, c.Custrole, c.NewAddresses})
}

type Notification struct {
	Delay    int    `json:"delay"`
	Schedule string `json:"schedule,omitempty"`
}

func (notification *Notification) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	notification.Delay = int(v["delay"].(float64))
	notification.Schedule = v["schedule"].(string)
	return nil
}

type CheckField struct {
	Name string `json:"name"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
}

type Schedule struct {
	Name       string                            `json:"id,omitempty"`
	CustomerId string                            `json:"customer_id,omitempty"`
	Data       map[string]map[string]interface{} `json:"data,omitempty"`
}

func (s *Schedule) MarshalJSONForCreate() ([]byte, error) {
	return json.Marshal(struct {
		Name       string                            `json:"id,omitempty"`
		CustomerId string                            `json:"customerid,omitempty"`
		Data       map[string]map[string]interface{} `json:"data,omitempty"`
	}{s.Name, s.CustomerId, s.Data})
}

type Group struct {
	ID         string   `json:"_id,omitempty"`
	CustomerId string   `json:"customer_id,omitempty"`
	Name       string   `json:"name"`
	Members    []string `json:"members"`
}

func (g *Group) MarshalJSONForCreate() ([]byte, error) {
	return json.Marshal(struct {
		CustomerId string   `json:"customerid,omitempty"`
		Name       string   `json:"name"`
		Members    []string `json:"members"`
	}{g.CustomerId, g.Name, g.Members})
}

type Customer struct {
	ID           string `json:"_id,omitempty"`
	Name         string `json:"customer_name,omitempty"`
	Email        string `json:"email,omitempty"`
	Parent       string `json:"parent,omitempty"`
	ContactName  string `json:"contact_name"`
	CreationDate int    `json:"creation_date,omitempty"`
	Status       string `json:"status"`
	Emailme      bool   `json:"emailme"`
	Timezone     string `json:"timezone"`
	Location     string `json:"location"`
}

func (c *Customer) MarshalJSONForCreate() ([]byte, error) {
	emailme := "no"
	if c.Emailme {
		emailme = "yes"
	}

	return json.Marshal(struct {
		Name        string `json:"name"`
		ContactName string `json:"contactname"`
		Email       string `json:"email"`
		Timezone    string `json:"timezone"`
		Location    string `json:"location"`
		Emailme     string `json:"emailme,omitempty"`
		Status      string `json:"status,omitempty"`
	}{c.Name, c.ContactName, c.Email, c.Timezone, c.Location, emailme, c.Status})
}

func (customer *Customer) UnmarshalJSON(data []byte) error {
	v := struct {
		ID           string   `json:"_id,omitempty"`
		Name         string   `json:"customer_name,omitempty"`
		Email        string   `json:"email,omitempty"`
		Parent       string   `json:"parent,omitempty"`
		ContactName  string   `json:"contact_name"`
		CreationDate int      `json:"creation_date,omitempty"`
		Status       string   `json:"status"`
		Emailme      bool     `json:"emailme"`
		Timezone     string   `json:"timezone"`
		Locations    []string `json:"defaultlocations"`
	}{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	customer.ID = v.ID
	customer.Name = v.Name
	customer.Email = v.Email
	customer.Parent = v.Parent
	customer.ContactName = v.ContactName
	customer.CreationDate = v.CreationDate
	customer.Status = v.Status
	customer.Emailme = v.Emailme
	customer.Timezone = v.Timezone
	customer.Location = v.Locations[0]
	return nil
}
