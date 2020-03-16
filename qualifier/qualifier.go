// Package qualifier - allows you to determine the type of Active Directory name format
package qualifier

import "unicode"

// Format - Active Directory name format
type Format int

func (f Format) String() string {
	switch f {
	case LDN:
		return "LegacyDomainName"
	case FQDN:
		return "Fully Qualified Domain Name"
	case UPN:
		return "User Principal Name"
	case DN:
		return "Distinguished Name"
	case CNAME:
		return "Canonical Name"
	case Display:
		return "Display"
	case GUID:
		return "GUID"
	case SID:
		return "SID"
	}
	return "Unknown"
}

const (
	// LDN - LegacyDomainName, which is also commonly referred to as the NetBIOS Domain Name
	LDN Format = iota // DOMAIN\username
	// UPN - User Principal Name
	UPN // username@domainname
	// DN - Distinguished Name
	DN // cn=JDoe Bond,ou=Widgets,ou=Manufacturing,dc=USRegion,dcOrgName.dc=com
	// GUID ...
	GUID // {95ee9fff-3436-11d1-b2b0-d15ae3ac8436}
	// CNAME - Canonical Name
	CNAME // USRegion.OrgName.com/Manufacturing/Widgets/JDoe Bond
	// FQDN - Fully Qualified Domain Name
	FQDN // subdomain1.subdomain2.domain.com
	// Display - display name
	Display // JDoe Bond
	// SID ...
	SID // S-1-5-21-1180699209-877415012-3182924384-1004
	// UNKNOWN - other format
	UNKNOWN
)

const (
	backslash    = '\\'
	forwardslash = '/'
	dot          = '.'
	equal        = '='
	dash         = '-'
	openbrace    = '{'
	lowerS       = 's'
	upperS       = 'S'
	at           = '@'
)

// Determine - to determine the type of Active Directory name format
func Determine(name string) Format {
	var stack rune
	var counter int
	for _, rune := range name {
		if unicode.IsSpace(rune) {
			continue
		}
		counter++
		switch rune {
		case openbrace:
			return GUID
		case backslash:
			return LDN
		case equal:
			return DN
		case at:
			return UPN
		case dot:
			stack = dot
		case forwardslash:
			return CNAME
		case dash:
			return SID
		}
	}
	if (stack == dot) && (counter > 1) {
		return FQDN
	}
	if counter > 0 {
		return Display
	}
	return UNKNOWN
}
