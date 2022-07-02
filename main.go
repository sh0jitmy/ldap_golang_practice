package main

import (
    //"crypto/tls"
    "errors"
    "log"
    "fmt"
    "github.com/go-ldap/ldap"
)

var (
    baseDN     = "dc=radio"
    username = "taro.kokusai"
    password = "5931"
    bindusername = "cn=admin,dc=radio"
    bindpassword = "password"
    ldapServer   = "ldap://localhost:30389"
)


func ExampleconnSearch() (bool, error) {

    err := errors.New("connection error")

    l, err := ldap.DialURL(ldapServer)
    if err != nil {
        fmt.Printf("%s\n", "ldap connection error")
        return false, err
    }

    fmt.Printf("%s\n", "ldap connection success")
    defer l.Close()
    return true, nil
}

func Example_userAuthentication() (bool,error) {
    err := errors.New("connection error")

    l, err := ldap.DialURL(ldapServer)
    if err != nil {
        fmt.Printf("%s\n", "ldap connection error")
        return false, err
    }

    fmt.Printf("%s\n", "ldap connection success")

    err = l.Bind(bindusername,bindpassword)
    if err != nil {
        fmt.Printf("%s\n", "ldap bind error")
        return false, err
    }

    // Search for the given username
    searchRequest := ldap.NewSearchRequest(
        baseDN,
        ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
        //fmt.Sprintf("(&(objectClass=organizationalPerson)(cn=%s))", username),
        fmt.Sprintf("(&(objectClass=organizationalPerson)\n"),
        //fmt.Sprintf("(&(cn=%s)\n",username),
        []string{"dn"},
        nil,
    )

    fmt.Printf("ldap search request\n")
    sr, err := l.Search(searchRequest)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("ldap search request end\n")

    if len(sr.Entries) != 1 {
        log.Fatal("User does not exist or too many entries returned")
    }

    userdn := sr.Entries[0].DN

    // Bind as the user to verify their password
    err = l.Bind(userdn, password)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("bind ok\n")

    // Rebind as the read only user for any further queries
    err = l.Bind(bindusername, bindpassword)
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()
    return true, nil
}


func main() {
    //ExampleconnSearch()
    Example_userAuthentication() 
}

