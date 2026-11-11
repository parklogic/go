package publicsuffix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWildcard(t *testing.T) {
	tests := []struct {
		name   string
		domain DomainName
		want   string
	}{
		{
			name: "no subdomain",
			domain: DomainName{
				FQDN:      "example.com",
				Domain:    "example.com",
				Subdomain: "",
				TLD:       suffixTLD,
			},
			want: "example.com",
		},
		{
			name: "single subdomain",
			domain: DomainName{
				FQDN:      "sub.example.com",
				Domain:    "example.com",
				Subdomain: "sub",
				TLD:       suffixTLD,
			},
			want: "*.example.com",
		},
		{
			name: "multiple subdomains",
			domain: DomainName{
				FQDN:      "sub1.sub2.example.com",
				Domain:    "example.com",
				Subdomain: "sub1.sub2",
				TLD:       suffixTLD,
			},
			want: "*.sub2.example.com",
		},
		{
			name: "nested subdomains",
			domain: DomainName{
				FQDN:      "sub1.sub2.sub3.example.com",
				Domain:    "example.com",
				Subdomain: "sub1.sub2.sub3",
				TLD:       suffixTLD,
			},
			want: "*.sub2.sub3.example.com",
		},
		{
			name: "edge case empty FQDN",
			domain: DomainName{
				FQDN:      "",
				Domain:    "",
				Subdomain: "",
				TLD:       "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.domain.Wildcard(); got != tt.want {
				t.Errorf("Wildcard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	name := DomainName{FQDN: "www.example.com"}
	assert.Equal(t, "www.example.com", name.String())
}
