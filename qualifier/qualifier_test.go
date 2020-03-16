// Package qualifier - allows you to determine the type of Active Directory name format

package qualifier

import (
	"reflect"
	"testing"
)

func TestDetermine(t *testing.T) {
	tests := []struct {
		name string
		want Format
	}{
		{" {95ee9fff-3436-11d1-b2b0-d15ae3ac8436} ", GUID},
		{"S-1-5-21-1180699209-877415012-3182924384-1004", SID},
		{"cn=JDoe Bond,ou=Widgets,ou=Manufacturing,dc=USRegion,dcOrgName.dc=com", DN},
		{"USRegion.OrgName.com/Manufacturing/Widgets/JDoe Bond", CNAME},
		{"J.Bond@domainname@USRegion.OrgName.com", UPN},
		{"USRegion\\J.Bond", LDN},
		{"  SJDoe Bond  ", Display},
		{" J.Bond.USRegion.OrgName.com ", FQDN},
		{"", UNKNOWN},
		{"    ", UNKNOWN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Determine(tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`Determine("%s") = %v, want %v`, tt.name, got, tt.want)
			}
		})
	}
}

func TestFormatString(t *testing.T) {
	tests := []struct {
		f    Format
		want string
	}{
		{LDN, "LegacyDomainName"},
		{FQDN, "Fully Qualified Domain Name"},
		{UPN, "User Principal Name"},
		{DN, "Distinguished Name"},
		{CNAME, "Canonical Name"},
		{Display, "Display"},
		{GUID, "GUID"},
		{SID, "SID"},
		{100, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("Format.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
