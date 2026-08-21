// Package abacpdp 实现 Attribute-Based Access Control 策略判决点（PDP）。
//
// 装载 PolicySet，绑定 Subject / Resource / Action / Environment 属性，
// 经属性匹配与合并算法（DenyOverrides 等）得到 Permit/Deny 与 Obligations。
// 本包专注判决，不含 HTTP 反向代理、功能开关投放或票据签发。
package abacpdp
