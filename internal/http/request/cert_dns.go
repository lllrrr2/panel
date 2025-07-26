package request

import "github.com/tnborg/panel/pkg/acme"

type CertDNSCreate struct {
	Type acme.DnsType  `form:"type" json:"type" validate:"required|in:aliyun,tencent,huawei,westcn,cloudflare,godaddy,gcore,porkbun,namecheap,namesilo,namecom,cloudns,duckdns,hetzner,linode,vercel"`
	Name string        `form:"name" json:"name" validate:"required"`
	Data acme.DNSParam `form:"data" json:"data" validate:"required"`
}

type CertDNSUpdate struct {
	ID   uint          `form:"id" json:"id" validate:"required|exists:cert_dns,id"`
	Type acme.DnsType  `form:"type" json:"type" validate:"required|in:aliyun,tencent,huawei,westcn,cloudflare,godaddy,gcore,porkbun,namecheap,namesilo,namecom,cloudns,duckdns,hetzner,linode,vercel"`
	Name string        `form:"name" json:"name" validate:"required"`
	Data acme.DNSParam `form:"data" json:"data" validate:"required"`
}
