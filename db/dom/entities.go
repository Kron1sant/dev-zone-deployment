package dom

// Represent an authenticated service user
type User struct {
	Id            uint   `bson:"_id" yaml:"id" json:"id"`
	Username      string `bson:"username" yaml:"username" json:"username"`
	EMail         string `bson:"email" yaml:"email" json:"email"`
	Password      string `bson:"password" yaml:"password" json:"-"`
	IsAdmin       bool   `bson:"is_admin" yaml:"is_admin" json:"isAdmin"`
	HasDevAccount bool   `bson:"hasDevAccount" yaml:"hasDevAccount" json:"hasDevAccount"`
	DevAccountId  uint   `bson:"devAccountId" yaml:"devAccountId" json:"devAccountId"`
}

// Represent a developer account
type DevAccount struct {
	Id             uint   `bson:"_id" json:"id"`
	Name           string `bson:"name" json:"name"`
	Surname        string `bson:"surname" json:"surname"`
	Patronomic     string `bson:"patronomic" json:"patronomic"`
	EMail          string `bson:"email" json:"email"`
	Username       string `bson:"username" json:"username"`
	HasOVPNCert    bool   `bson:"hasOVPNCert" json:"hasOVPNCert"`
	OpenVPNKeyName string `bson:"openVPNKeyName" json:"-"`
	Comment        string `bson:"comment" json:"comment"`
}

// Represent an instance of Virtual machine
type VM struct {
	Id            string   `bson:"_id" json:"id"`
	Name          string   `bson:"name" json:"name"`
	HasDevAccount bool     `bson:"hasDevAccount" yaml:"hasDevAccount" json:"hasDevAccount"`
	DevAccountId  uint     `bson:"devAccountId" yaml:"devAccountId" json:"devAccountId"`
	Description   string   `bson:"description" json:"description"`
	Params        string   `bson:"params" json:"params"`
	Status        StatusVM `bson:"status" json:"status"`
}

type StatusVM string

const (
	VM_STATUS_PROVISIONING StatusVM = "PROVISIONING"
	VM_STATUS_STARTING     StatusVM = "STARTING"
	VM_STATUS_RUNNING      StatusVM = "RUNNING"
	VM_STATUS_STOPPING     StatusVM = "STOPPING"
	VM_STATUS_STOPPED      StatusVM = "STOPPED"
	VM_STATUS_RESTARTING   StatusVM = "RESTARTING"
	VM_STATUS_UPDATING     StatusVM = "UPDATING"
	VM_STATUS_CRASHED      StatusVM = "CRASHED"
	VM_STATUS_ERROR        StatusVM = "ERROR"
	VM_STATUS_DELETING     StatusVM = "DELETING"
)
