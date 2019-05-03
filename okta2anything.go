package main  //hi

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	ldap "github.com/vjeantet/ldapserver"
	"strings"
	message "github.com/vjeantet/goldap/message"
	"time"
	"math/rand"
	"reflect"
	"flag"
	"os/exec"
	"encoding/json"
)

//Structs
type LdapStruct struct {
	Active           	       string `json:"active"`
	DepartmentNumber           string `json:"departmentNumber"`
	Distinguishedname          string `json:"distinguishedname"`
	Ds_pwp_account_disabled    string `json:"ds-pwp-account-disabled"`
	EmployeeID                 string `json:"employeeID"`
	EmployeeNumber             string `json:"employeeNumber"`
	Entryuuid                  string `json:"entryuuid"`
	GivenName                  string `json:"givenName"`
	L                          string `json:"l"`
	Mail                       string `json:"mail"`
	Manager                    string `json:"manager"`
	Mobile                     string `json:"mobile"`
	PhysicalDeliveryOfficeName string `json:"physicalDeliveryOfficeName"`
	PostalAddress              string `json:"postalAddress"`
	PostalCode                 string `json:"postalCode"`
	PreferredLanguage          string `json:"preferredLanguage"`
	Sn                         string `json:"sn"`
	St                         string `json:"st"`
	StreetAddress              string `json:"streetAddress"`
	SupportedLDAPVersion       string `json:"supportedLDAPVersion"`
	TelephoneNumber            string `json:"telephoneNumber"`
	Title                      string `json:"title"`
	UID                        string `json:"uid"`
}
//End Structs

//Globals

var (
	//LdapBind        string
	LdapPassword        string
	Plugin        string
	PostPlugin        string
	PrePlugin        string
	verboseOutput bool
)

var stackMap = make(map[string]string) // map with username as index



//end Globals

func init() {


	verboseOutput = false

	rand.Seed(time.Now().UnixNano())  //Seed the randomizer

	//Check Commnad line arguments

	//binddn := flag.String("D", "Promiscuous Mode", "Directory Manager Example: -D=\"Directory Manager\"")
	password := flag.String("w", "Promiscuous Mode", "Password for Directory Manager Example: -w=Password1")
	plugin := flag.String("plugin", "", "External Authorizer")
	postplugin := flag.String("postplugin", "", "Post Authentication function")
	preplugin := flag.String("preplugin", "", "Pre Authentication function")
	flag.Parse()

	//LdapBind =*binddn
	LdapPassword=*password
	Plugin=*plugin
	PostPlugin=*postplugin
	PrePlugin=*preplugin



	fmt.Println("Okta2Anything, For more command line options use the --help switch")
	fmt.Println("               Directory Manager set to cn=Directory Manager")

	if (len(Plugin)==0) {
		fmt.Println("\n ERROR, you need to specify a Plugin, use -plugin=promiscuous for testing")
		os.Exit(0)
	}
	//if (LdapBind=="Promiscuous Mode") {
	//	fmt.Println("Running in Promiscous Mode! All Authentications are permitted\n")
	//
	//}
	//End Check Commandline arguments


}

func main() {

	//Create a new LDAP Server
	server := ldap.NewServer()

	//Create routes bindings
	routes := ldap.NewRouteMux()
	routes.NotFound(handleNotFound)
	routes.Abandon(handleAbandon)
	routes.Bind(handleBind)
	routes.Compare(handleCompare)
	routes.Add(handleAdd)
	routes.Delete(handleDelete)
	routes.Modify(handleModify)
	/* */

	routes.Extended(handleStartTLS).
		RequestName(ldap.NoticeOfStartTLS).Label("StartTLS")

	routes.Extended(handleWhoAmI).
		RequestName(ldap.NoticeOfWhoAmI).Label("Ext - WhoAmI")

	routes.Extended(handleExtended).Label("Ext - Generic")

	routes.Search(handleSearchDSE).
		BaseDn("").
		Scope(ldap.SearchRequestScopeBaseObject).
		Filter("(objectclass=*)").
		Label("Search - ROOT DSE")

	routes.Search(handleSearchMyCompany).
		BaseDn("ou=People,dc=example,dc=com").
		Scope(ldap.SearchRequestScopeBaseObject).
		Label("Search - Compagny Root")

	routes.Search(handleSearch).Label("Search - Generic")

	//Attach routes to server
	server.Handle(routes)

	// listen on 10389 and serve
	go server.ListenAndServe("0.0.0.0:389")

	// When CTRL+C, SIGINT and SIGTERM signal occurs
	// Then stop server gracefully
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)

	server.Stop()
}

func handleNotFound(w ldap.ResponseWriter, r *ldap.Message) {
	switch r.ProtocolOpType() {
	case ldap.ApplicationBindRequest:
		res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
		res.SetDiagnosticMessage("Default binding behavior set to return Success")

		w.Write(res)

	default:
		res := ldap.NewResponse(ldap.LDAPResultUnwillingToPerform)
		res.SetDiagnosticMessage("Operation not implemented by server")
		w.Write(res)
	}
}

func handleAbandon(w ldap.ResponseWriter, m *ldap.Message) {
	var req = m.GetAbandonRequest()
	// retreive the request to abandon, and send a abort signal to it
	if requestToAbandon, ok := m.Client.GetMessageByID(int(req)); ok {
		requestToAbandon.Abandon()
		log.Printf("Abandon signal sent to request processor [messageID=%d]", int(req))
	}
}

func handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	log.Printf("HIT handleBind Function #188 !\n")

	r := m.GetBindRequest()

	//log.Printf("Bind Attempt User=%s, Pass=XXXXXXXXX", string(r.Name()), r.Authentication())
	log.Printf("Bind Attempt User=%s, Pass=XXXXXXXXX", string(r.Name()))
	password:=fmt.Sprintf("%s",r.Authentication())


	//r := m.GetBindRequest()
	//res := ldap.NewBindResponse(ldap.LDAPResultSuccess)



	if r.AuthenticationChoice() == "simple" {


		// Commented out for promiscuis
		if string(r.Name()) == "cn=Directory Manager" {
			log.Printf("Hello Directory Manager\n")
			if (LdapPassword=="Promiscuous Mode") {
				log.Printf("Directory Manager running in Promiscuous Mode\n")

				res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			} else if (LdapPassword==password) {
				log.Printf("Directory Manager password is correct\n")

				res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			} else {
				log.Printf("Directory Manager Passwords do not Match!\n")
				res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
				w.Write(res)
				return
			}

		}

		if len(Plugin) !=0 {

			ldapObj :=LdapStruct{}

			out, err := exec.Command("node", "./"+Plugin,string(r.Name()),password).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("--------------------------------------")
			fmt.Println("230: --------------------------------------")
			fmt.Println(string(out));
			pushToStack(string(r.Name()), string(out))
			fmt.Println("I think name is"+r.Name())
			eraseme := popFromStack(string(r.Name()))
			fmt.Println(eraseme)
			json.Unmarshal(out,&ldapObj)


			//fmt.Println(ldapObj.Active)
			if ldapObj.Active=="true" {
				res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
				w.Write(res)

				//add post plugin

				if len(PostPlugin) !=0 {

					fmt.Println ("Perform Post Plugin for :"+r.Name());

					out2, err2 := exec.Command("node", "./"+PostPlugin,string(r.Name()),password).Output()
					if err2 != nil {
						log.Fatal(err)
					}
					_=out2
				}

				//end post Plugin

				return
			} else {
				res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
				w.Write(res)
				return
			}

		}

		//if LdapBind=="Promiscuous Mode" {
		//	res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
		//	w.Write(res)
		//	return
		//}
	}

	zp:=fmt.Sprintf("%s",r.Authentication())

	if ( zp == "Password1") {
		res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		return

	} else {
		res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
		w.Write(res)
		return
	}

	//temp hack all auth
	//res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
	//
	//w.Write(res)
	//	return
	// end temp hack
	//	log.Printf("Bind failed User=%s, Pass=%#v", string(r.Name()), r.Authentication())
	//	log.Printf("#124")
	//	res.SetResultCode(ldap.LDAPResultInvalidCredentials)
	//	res.SetDiagnosticMessage("invalid credentials")
	//} else {
	//	res.SetResultCode(ldap.LDAPResultUnwillingToPerform) //Really unwilling.. Funny
	//	res.SetDiagnosticMessage("Authentication choice not supported")
	//}
	//
	//w.Write(res)
}

// The resultCode is set to compareTrue, compareFalse, or an appropriate
// error.  compareTrue indicates that the assertion value in the ava
// Comparerequest field matches a value of the attribute or subtype according to the
// attribute's EQUALITY matching rule.  compareFalse indicates that the
// assertion value in the ava field and the values of the attribute or
// subtype did not match.  Other result codes indicate either that the
// result of the comparison was Undefined, or that
// some error occurred.
func handleCompare(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetCompareRequest()
	log.Printf("Comparing entry: %s", r.Entry())
	//attributes values
	log.Printf(" attribute name to compare : \"%s\"", r.Ava().AttributeDesc())
	log.Printf(" attribute value expected : \"%s\"", r.Ava().AssertionValue())

	res := ldap.NewCompareResponse(ldap.LDAPResultCompareTrue)

	w.Write(res)
}

func handleAdd(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetAddRequest()
	log.Printf("Adding entry: %s", r.Entry())
	//attributes values
	for _, attribute := range r.Attributes() {
		for _, attributeValue := range attribute.Vals() {
			log.Printf("- %s:%s", attribute.Type_(), attributeValue)
		}
	}
	res := ldap.NewAddResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleModify(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetModifyRequest()
	log.Printf("Modify entry: %s", r.Object())

	for _, change := range r.Changes() {
		modification := change.Modification()
		var operationString string
		switch change.Operation() {
		case ldap.ModifyRequestChangeOperationAdd:
			operationString = "Add"
		case ldap.ModifyRequestChangeOperationDelete:
			operationString = "Delete"
		case ldap.ModifyRequestChangeOperationReplace:
			operationString = "Replace"
		}

		log.Printf("%s attribute '%s'", operationString, modification.Type_())
		for _, attributeValue := range modification.Vals() {
			log.Printf("- value: %s", attributeValue)
		}

	}

	res := ldap.NewModifyResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleDelete(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetDeleteRequest()
	log.Printf("Deleting entry: %s", r)
	res := ldap.NewDeleteResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleExtended(w ldap.ResponseWriter, m *ldap.Message) {
	log.Printf("HIT handleExtended Function !\n")

	r := m.GetExtendedRequest()
	log.Printf("Extended request received, name=%s", r.RequestName())
	log.Printf("Extended request received, value=%x", r.RequestValue())
	res := ldap.NewExtendedResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleWhoAmI(w ldap.ResponseWriter, m *ldap.Message) {
	res := ldap.NewExtendedResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleSearchDSE(w ldap.ResponseWriter, m *ldap.Message) {
	log.Printf("HIT handleSearchDSE Function #216 !\n")

	r := m.GetSearchRequest()



	log.Printf("Search DSE")
	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())
	log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())

	e := ldap.NewSearchResultEntry("")
	e.AddAttribute("vendorName", "Valère JEANTET")
	e.AddAttribute("vendorVersion", "0.0.1")
	e.AddAttribute("objectClass", "top", "extensibleObject")
	e.AddAttribute("supportedLDAPVersion", "3")
	e.AddAttribute("namingContexts", "dc=example,dc=com")
	e.AddAttribute("departmentNumber", "11")
	e.AddAttribute("telephoneNumber", "1111111")
	// e.AddAttribute("subschemaSubentry", "cn=schema")
	// e.AddAttribute("namingContexts", "ou=system", "ou=schema", "dc=example,dc=com", "ou=config")
	// e.AddAttribute("supportedFeatures", "1.3.6.1.4.1.4203.1.5.1")
	// e.AddAttribute("supportedControl", "2.16.840.1.113730.3.4.3", "1.3.6.1.4.1.4203.1.10.1", "2.16.840.1.113730.3.4.2", "1.3.6.1.4.1.4203.1.9.1.4", "1.3.6.1.4.1.42.2.27.8.5.1", "1.3.6.1.4.1.4203.1.9.1.1", "1.3.6.1.4.1.4203.1.9.1.3", "1.3.6.1.4.1.4203.1.9.1.2", "1.3.6.1.4.1.18060.0.0.1", "2.16.840.1.113730.3.4.7", "1.2.840.113556.1.4.319")
	// e.AddAttribute("supportedExtension", "1.3.6.1.4.1.1466.20036", "1.3.6.1.4.1.4203.1.11.1", "1.3.6.1.4.1.18060.0.1.5", "1.3.6.1.4.1.18060.0.1.3", "1.3.6.1.4.1.1466.20037")
	// e.AddAttribute("supportedSASLMechanisms", "NTLM", "GSSAPI", "GSS-SPNEGO", "CRAM-MD5", "SIMPLE", "DIGEST-MD5")
	// e.AddAttribute("entryUUID", "f290425c-8272-4e62-8a67-92b06f38dbf5")
	w.Write(e)

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
	fmt.Println("Exiting 419 !!!!!!!!!!!!!!!!\n\n")
}

func handleSearchMyCompany(w ldap.ResponseWriter, m *ldap.Message) {
	log.Printf("HIT handleSearchMyCompany Function !\n")

	r := m.GetSearchRequest()
	log.Printf("handleSearchMyCompany - Request BaseDn=%s", r.BaseObject())

	e := ldap.NewSearchResultEntry(string(r.BaseObject()))
	e.AddAttribute("objectClass", "top", "organizationalUnit")
	w.Write(e)

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
	log.Printf("HIT handleSearch Function #380 !\n")

	r := m.GetSearchRequest()

	userid:="aaaaa"

	log.Printf("Search DSE")
	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())
	log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())


	getUuid:=fmt.Sprintf("%s",r.FilterString())

	if ( strings.Index(getUuid,"entryuuid=") !=0) {
		getUuid = strings.Replace(getUuid, "(", "", -1)
		getUuid = strings.Replace(getUuid, ")", "", -1)
		getUuid = strings.Replace(getUuid, "&", "", -1)

		entryuuidIndex := strings.Index(getUuid, "entryuuid=")
		getUuid = getUuid[entryuuidIndex+10:]

		userid=getUuid
	} else {
		fmt.Println("463: weird condition")
		getUuid="xxxxxxx"
	}

	// ou=groups search .. Treat Groups differently
	if strings.HasPrefix(string(r.BaseObject()),"ou=Groups") ||
		strings.HasPrefix(string(r.BaseObject()),"ou=groups") {

		log.Printf("Search DSE Groups #402")


		select {
		case <-m.Done:
			log.Print("Leaving handleSearch...")
			return
		default:
		}

		e := ldap.NewSearchResultEntry("ou=Groups,dc=example,dc=com")

		e.AddAttribute("objectClass", "top", "organizationalUnit")
		e.AddAttribute("ou", "Groups")
		e.AddAttribute(message.AttributeDescription( "entryuuid"), message.AttributeValue(getUuid))


		w.Write(e)

		e = ldap.NewSearchResultEntry("cn=ldapusers,ou=Groups,dc=example,dc=com")

		e.AddAttribute("objectClass", "top", "groupOfNames")
		e.AddAttribute("cn", "ldapusers")

		e.AddAttribute(message.AttributeDescription("member"), message.AttributeValue("uid"+userid+"=user.0,ou=People,dc=example,dc=com"))

		e.AddAttribute("entryuuid", "8c3624f5-d219-4401-9042-9a1fbf6f1b6805")
		w.Write(e)

		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		fmt.Println("Exiting 500!!!! Group Search!\n\n")

	} else {

		fmt.Println("506: second things after groups")

		// Handle Stop Signal (server stop / client disconnected / Abandoned request....)
		select {
		case <-m.Done:
			log.Print("Leaving handleSearch...")
			fmt.Println("512: Might have blown up")
			return
		default:
		}

		log.Printf("FileterString: %s",reflect.TypeOf(r.FilterString()))

		FilterString:= r.FilterString()
		if strings.Contains(FilterString,"(&(objectclass=inetorgperson)(mail=") {
			fmt.Println("521 trying to figure out username")
			FilterString=strings.Replace(FilterString,"(&(objectclass=inetorgperson)(mail=", "", -1)
			FilterString=strings.Replace(FilterString,")", "", -1)
			log.Printf("FilterString: %s",FilterString)
			log.Printf("Line #325")
			userid = FilterString
		}

		e := ldap.NewSearchResultEntry("uid="+userid+",ou=People,dc=example,dc=com")

		e.AddAttribute("objectClass", "top", "inetorgperson", "organizationalPerson", "person")

		if userid == "ss=person" {
			fmt.Println("trying to get userid again")
			fmt.Println(r.BaseObject());
			userid=fmt.Sprintf("%s",r.BaseObject())
			fmt.Println("569: New userid = "+userid)

		}

		if strings.Contains(userid, "@") {

		} else {
			userid=userid+"@noemailprovided.com"
		}

		fmt.Println("Here is the userid:"+userid)

		e.AddAttribute(message.AttributeDescription("distinguishedname"), message.AttributeValue(userid))
		e.AddAttribute(message.AttributeDescription("mail"), message.AttributeValue(userid))

//		e.AddAttribute(message.AttributeDescription("givenName"), message.AttributeValue(userid))
//		e.AddAttribute(message.AttributeDescription("sn"), message.AttributeValue(userid))

		e.AddAttribute("supportedLDAPVersion", "3")
		e.AddAttribute("title", "title")
		e.AddAttribute(message.AttributeDescription("uid"), message.AttributeValue(userid))
		e.AddAttribute("manager", "manager")
		e.AddAttribute("streetAddress", "street")
		e.AddAttribute("l", "USA")
		e.AddAttribute("st", "TX")

		// Add custom attributes

		var customAttributes=popFromStack(userid)

		fmt.Println("===============================")
		fmt.Println("I'm adding")
		fmt.Println("===============================")
		fmt.Println(customAttributes)
		fmt.Println("END adding custom")

		e.AddAttribute(message.AttributeDescription("givenName"), message.AttributeValue(userid))
		e.AddAttribute(message.AttributeDescription("sn"), message.AttributeValue(userid))

		fmt.Println(len(customAttributes))

		fmt.Println(PrePlugin);

		if len(PrePlugin) !=0 {


		if len(customAttributes) == 0 {
			fmt.Println("Map is Empty")
			//e.AddAttribute(message.AttributeDescription("departmentNumber"), message.AttributeValue(userid))
			//e.AddAttribute(message.AttributeDescription("telephoneNumber"), message.AttributeValue(userid))

			preout, err := exec.Command("node", "./"+PrePlugin, string("1"), "2").Output()
			if err != nil {
				fmt.Println("Error")
				log.Fatal(err)
			}

			fmt.Println(string(preout))

			preloginFields := convertJsonStringToMap(string(preout))

			fmt.Println(preloginFields)

			for customkey, customvalue := range preloginFields {
				fmt.Println("598:  Adding prelogin fields")
				fmt.Printf("%s -> %s\n", customkey, customvalue)
				e.AddAttribute(message.AttributeDescription(fmt.Sprintf("%s", string(customkey))),
					message.AttributeValue(fmt.Sprintf("%s", string(customvalue))))

			}
		}

		} else {

			for customkey, customvalue := range customAttributes {
				fmt.Println("I'm in the Loop !!")
				fmt.Printf("%s -> %s\n", customkey, customvalue)
				e.AddAttribute(message.AttributeDescription(fmt.Sprintf("%s", string(customkey))),
					message.AttributeValue(fmt.Sprintf("%s", string(customvalue))))

			}
		}

		e.AddAttribute(message.AttributeDescription(fmt.Sprintf("%s", string ("pizza"))),
			message.AttributeValue(fmt.Sprintf("%s", string ("guidguidhere"))))

		e.AddAttribute("postalCode", "11111")
		e.AddAttribute("physicalDeliveryOfficeName", "x")
		//e.AddAttribute("departmentNumber", "11")
		//e.AddAttribute("telephoneNumber", "1111111")
		e.AddAttribute("mobile", "1111111")
		e.AddAttribute("preferredLanguage", "en")
		e.AddAttribute("postalAddress", "Austin")
		e.AddAttribute("employeeID", "0")
		e.AddAttribute("employeeNumber", "0")
		e.AddAttribute(message.AttributeDescription("entryuuid"),
			message.AttributeValue(string (userid)))

		e.AddAttribute("ds-pwp-account-disabled", "")
		w.Write(e)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		fmt.Println("Exiting 596 Done !! ---------------\n\n")
	}

}

// localhostCert is a PEM-encoded TLS cert with SAN DNS names
// "127.0.0.1" and "[::1]", expiring at the last second of 2049 (the end
// of ASN.1 time).
var localhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIBOTCB5qADAgECAgEAMAsGCSqGSIb3DQEBBTAAMB4XDTcwMDEwMTAwMDAwMFoX
DTQ5MTIzMTIzNTk1OVowADBaMAsGCSqGSIb3DQEBAQNLADBIAkEAsuA5mAFMj6Q7
qoBzcvKzIq4kzuT5epSp2AkcQfyBHm7K13Ws7u+0b5Vb9gqTf5cAiIKcrtrXVqkL
8i1UQF6AzwIDAQABo08wTTAOBgNVHQ8BAf8EBAMCACQwDQYDVR0OBAYEBAECAwQw
DwYDVR0jBAgwBoAEAQIDBDAbBgNVHREEFDASggkxMjcuMC4wLjGCBVs6OjFdMAsG
CSqGSIb3DQEBBQNBAJH30zjLWRztrWpOCgJL8RQWLaKzhK79pVhAx6q/3NrF16C7
+l1BRZstTwIGdoGId8BRpErK1TXkniFb95ZMynM=
-----END CERTIFICATE-----
`)

// localhostKey is the private key for localhostCert.
var localhostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBPQIBAAJBALLgOZgBTI+kO6qAc3LysyKuJM7k+XqUqdgJHEH8gR5uytd1rO7v
tG+VW/YKk3+XAIiCnK7a11apC/ItVEBegM8CAwEAAQJBAI5sxq7naeR9ahyqRkJi
SIv2iMxLuPEHaezf5CYOPWjSjBPyVhyRevkhtqEjF/WkgL7C2nWpYHsUcBDBQVF0
3KECIQDtEGB2ulnkZAahl3WuJziXGLB+p8Wgx7wzSM6bHu1c6QIhAMEp++CaS+SJ
/TrU0zwY/fW4SvQeb49BPZUF3oqR8Xz3AiEA1rAJHBzBgdOQKdE3ksMUPcnvNJSN
poCcELmz2clVXtkCIQCLytuLV38XHToTipR4yMl6O+6arzAjZ56uq7m7ZRV0TwIh
AM65XAOw8Dsg9Kq78aYXiOEDc5DL0sbFUu/SlmRcCg93
-----END RSA PRIVATE KEY-----
`)

// getTLSconfig returns a tls configuration used
// to build a TLSlistener for TLS or StartTLS
func getTLSconfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair(localhostCert, localhostKey)
	if err != nil {
		return &tls.Config{}, err
	}

	return &tls.Config{
		MinVersion:   tls.VersionSSL30,
		MaxVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ServerName:   "127.0.0.1",
	}, nil
}

func handleStartTLS(w ldap.ResponseWriter, m *ldap.Message) {
	tlsconfig, _ := getTLSconfig()
	tlsConn := tls.Server(m.Client.GetConn(), tlsconfig)
	res := ldap.NewExtendedResponse(ldap.LDAPResultSuccess)
	res.SetResponseName(ldap.NoticeOfStartTLS)
	w.Write(res)

	if err := tlsConn.Handshake(); err != nil {
		log.Printf("StartTLS Handshake error %v", err)
		res.SetDiagnosticMessage(fmt.Sprintf("StartTLS Handshake error : \"%s\"", err.Error()))
		res.SetResultCode(ldap.LDAPResultOperationsError)
		w.Write(res)
		return
	}

	m.Client.SetConn(tlsConn)
	log.Println("StartTLS OK")
}


func RandStringR(n int) string {

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandInt(n int) string {

	var letterRunes = []rune("0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func pushToStack ( key string, jsonString string ) {
	fmt.Println("((((((((((((((((((((((((((((((((((((((((((")
	fmt.Println("pushToStack adding user:"+key)
	fmt.Println("pushToStack adding data:"+jsonString)
	fmt.Println("((((((((((((((((((((((((((((((((((((((((((")

	stackMap[key]=jsonString
}

func popFromStack ( key string) map [string]string {

	fmt.Println("Looking for:",key)
	fmt.Println(stackMap)

	if val, ok := stackMap[key]; ok { //Make ure there is a match
		//delete(stackMap, key) //remove it
		return convertJsonStringToMap(val)
	}

	return make (map[string]string)  //Didn't find your key.. Sorry about that
	// returning empty map
}

func convertJsonStringToMap ( jsonData string ) map [string]string {

	var mapToReturn=make (map[string]string)

	jsonByteArray:= []byte(jsonData)
	var v interface{}
	err:=json.Unmarshal(jsonByteArray, &v)
	if err!=nil {
		fmt.Println("****** JSON Parse Error *****\n", jsonData)
		return mapToReturn //Something Blew up parsing the JSON
	}
	//fmt.Println(err)
	data := v.(map[string]interface{})

	for k, v := range data {
		if (k!="Active") {
			valueToString, ok := v.(string)
			_ = ok
			mapToReturn[string(k)] = valueToString
		}
	}

	return mapToReturn



}



