package goreferrer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlankReferrerIsDirect(t *testing.T) {
	blank := DefaultRules.Parse("")
	whitespace := DefaultRules.Parse(" \t\n\r")
	assert.Equal(t, Direct, blank.Type)
	assert.Equal(t, Direct, whitespace.Type)
}

func TestUnknownDomainIsIndirect(t *testing.T) {
	actual := DefaultRules.Parse("https://example.org/path?foo=bar#baz")
	expected := Referrer{
		Type:   Indirect,
		URL:    "https://example.org/path?foo=bar#baz",
		Host:   "example.org",
		Domain: "example",
		Tld:    "org",
		Path:   "/path",
	}
	assert.Equal(t, expected, actual)
}

func TestMalformedUrlIsInvalid(t *testing.T) {
	urls := []string{
		"blap",
		"blap blap",
		"http://",
		"/",
	}

	for _, u := range urls {
		if !assert.Equal(t, Invalid, DefaultRules.Parse(u).Type) {
			t.Log(u)
		}
	}
}

func TestMatchOnAllButQueryString(t *testing.T) {
	rules := RuleSet{
		"www.zambo.com/search": Rule{Type: Search},
	}
	assert.Equal(t, Search, rules.Parse("http://www.zambo.com/search?q=hello!").Type)
}

func TestMatchOnDomainTldAndPath(t *testing.T) {
	rules := RuleSet{
		"zambo.com/search": Rule{Type: Search},
	}
	assert.Equal(t, Search, rules.Parse("http://www.zambo.com/search?q=hello!").Type)
}

func TestMatchOnSubdomainDomainAndTld(t *testing.T) {
	rules := RuleSet{
		"www.zambo.com": Rule{Type: Search},
	}
	assert.Equal(t, Search, rules.Parse("http://www.zambo.com/search?q=hello!").Type)
}

func TestMatchOnDomainAndTld(t *testing.T) {
	rules := RuleSet{
		"zambo.com": Rule{Type: Search},
	}
	assert.Equal(t, Search, rules.Parse("http://www.zambo.com/search?q=hello!").Type)
}

func TestSimpleEmail(t *testing.T) {
	actual := DefaultRules.Parse("https://mail.google.com/9aifaufasodf8usafd")
	expected := Referrer{
		Type:      Email,
		Label:     "Gmail",
		URL:       "https://mail.google.com/9aifaufasodf8usafd",
		Host:      "mail.google.com",
		Subdomain: "mail",
		Domain:    "google",
		Tld:       "com",
		Path:      "/9aifaufasodf8usafd",
	}
	assert.Equal(t, expected, actual)
}

func TestSimpleSocial(t *testing.T) {
	actual := DefaultRules.Parse("https://twitter.com/snormore/status/391149968360103936")
	expected := Referrer{
		Type:   Social,
		Label:  "Twitter",
		URL:    "https://twitter.com/snormore/status/391149968360103936",
		Host:   "twitter.com",
		Domain: "twitter",
		Tld:    "com",
		Path:   "/snormore/status/391149968360103936",
	}
	assert.Equal(t, expected, actual)
}

func TestSocialSubdomain(t *testing.T) {
	actual := DefaultRules.Parse("https://puppyanimalbarn.tumblr.com")
	expected := Referrer{
		Type:      Social,
		Label:     "Tumblr",
		URL:       "https://puppyanimalbarn.tumblr.com",
		Host:      "puppyanimalbarn.tumblr.com",
		Subdomain: "puppyanimalbarn",
		Domain:    "tumblr",
		Tld:       "com",
	}
	assert.Equal(t, expected, actual)
}
