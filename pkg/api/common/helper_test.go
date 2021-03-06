// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package common

import (
	"testing"

	"github.com/pkg/errors"
)

func TestValidateDNSPrefix(t *testing.T) {
	cases := []struct {
		dnsPrefix   string
		expectedErr error
	}{
		{
			"validDnsPrefix",
			nil,
		},
		{
			"",
			errors.New("DNSPrefix '' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 0)"),
		},
		{
			"a",
			errors.New("DNSPrefix 'a' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 1)"),
		},
		{
			"1234",
			errors.New("DNSPrefix '1234' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 4)"),
		},
		{
			"verylongdnsprefixthatismorethan45characterslong",
			errors.New("DNSPrefix 'verylongdnsprefixthatismorethan45characterslong' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 47)"),
		},
		{
			"dnswith_special?char",
			errors.New("DNSPrefix 'dnswith_special?char' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 20)"),
		},
		{
			"myDNS-1234",
			nil,
		},
	}

	for _, c := range cases {
		err := ValidateDNSPrefix(c.dnsPrefix)
		if err != nil && c.expectedErr != nil {
			if err.Error() != c.expectedErr.Error() {
				t.Fatalf("expected validateDNSPrefix to return error %s, but instead got %s", c.expectedErr.Error(), err.Error())
			}
		} else {
			if c.expectedErr != nil {
				t.Fatalf("expected validateDNSPrefix to return error %s, but instead got no error", c.expectedErr.Error())
			} else if err != nil {
				t.Fatalf("expected validateDNSPrefix to return no error, but instead got %s", err.Error())
			}
		}
	}
}

func TestIsNvidiaEnabledSKU(t *testing.T) {
	cases := GetNSeriesVMCasesForTesting()

	for _, c := range cases {
		ret := IsNvidiaEnabledSKU(c.VMSKU)
		if ret != c.Expected {
			t.Fatalf("expected IsNvidiaEnabledSKU(%s) to return %t, but instead got %t", c.VMSKU, c.Expected, ret)
		}
	}
}

func getCSeriesVMCasesForTesting() []struct {
	VMSKU    string
	Expected bool
} {
	cases := []struct {
		VMSKU    string
		Expected bool
	}{
		{
			"Standard_DC2s",
			true,
		},
		{
			"Standard_DC4s",
			true,
		},
		{
			"Standard_D2_v2",
			false,
		},
		{
			"gobledygook",
			false,
		},
		{
			"",
			false,
		},
	}
	return cases
}

func TestIsSGXEnabledSKU(t *testing.T) {
	cases := getCSeriesVMCasesForTesting()

	for _, c := range cases {
		ret := IsSgxEnabledSKU(c.VMSKU)
		if ret != c.Expected {
			t.Fatalf("expected IsSgxEnabledSKU(%s) to return %t, but instead got %t", c.VMSKU, c.Expected, ret)
		}
	}
}

func TestGetMasterKubernetesLabels(t *testing.T) {
	cases := []struct {
		rg       string
		expected string
	}{
		{
			"my-resource-group",
			"kubernetes.io/role=master,node-role.kubernetes.io/master=,kubernetes.azure.com/cluster=my-resource-group",
		},
		{
			"",
			"kubernetes.io/role=master,node-role.kubernetes.io/master=,kubernetes.azure.com/cluster=",
		},
	}

	for _, c := range cases {
		ret := GetMasterKubernetesLabels(c.rg)
		if ret != c.expected {
			t.Fatalf("expected GetMasterKubernetesLabels(%s) to return %s, but instead got %s", c.rg, c.expected, ret)
		}
	}
}

func TestGetOrderedEscapedKeyValsString(t *testing.T) {
	alphabetizedString := `\"foo=bar\", \"yes=please\"`
	cases := []struct {
		input    map[string]string
		expected string
	}{
		{
			input:    map[string]string{},
			expected: "",
		},
		{
			input: map[string]string{
				"foo": "bar",
				"yes": "please",
			},
			expected: alphabetizedString,
		},
		{
			input: map[string]string{
				"yes": "please",
				"foo": "bar",
			},
			expected: alphabetizedString,
		},
	}

	for _, c := range cases {
		ret := GetOrderedEscapedKeyValsString(c.input)
		if ret != c.expected {
			t.Fatalf("expected GetOrderedEscapedKeyValsString(%s) to return %s, but instead got %s", c.input, c.expected, ret)
		}
	}
}

func TestGetStorageAccountType(t *testing.T) {
	validPremiumVMSize := "Standard_DS2_v2"
	validStandardVMSize := "Standard_D2_v2"
	expectedPremiumTier := "Premium_LRS"
	expectedStandardTier := "Standard_LRS"
	invalidVMSize := "D2v2"

	// test premium VMSize returns premium managed disk tier
	premiumTier, err := GetStorageAccountType(validPremiumVMSize)
	if err != nil {
		t.Fatalf("Invalid sizeName: %s", err)
	}

	if premiumTier != expectedPremiumTier {
		t.Fatalf("premium VM did no match premium managed storage tier")
	}

	// test standard VMSize returns standard managed disk tier
	standardTier, err := GetStorageAccountType(validStandardVMSize)
	if err != nil {
		t.Fatalf("Invalid sizeName: %s", err)
	}

	if standardTier != expectedStandardTier {
		t.Fatalf("standard VM did no match standard managed storage tier")
	}

	// test invalid VMSize
	result, err := GetStorageAccountType(invalidVMSize)
	if err == nil {
		t.Errorf("GetStorageAccountType() = (%s, nil), want error", result)
	}
}

func TestSliceIntIsNonEmpty(t *testing.T) {
	cases := []struct {
		input    []int
		expected bool
	}{
		{
			input: []int{
				1, 2, 3,
			},
			expected: true,
		},
		{
			input:    []int{},
			expected: false,
		},
		{
			input:    nil,
			expected: false,
		},
	}

	for _, c := range cases {
		ret := SliceIntIsNonEmpty(c.input)
		if ret != c.expected {
			t.Fatalf("expected SliceIntIsNonEmpty(%v) to return %t, but instead got %t", c.input, c.expected, ret)
		}
	}
}
