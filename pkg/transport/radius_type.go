/*
 Copyright © 2024 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package transport

import (
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"layeh.com/radius/rfc2867"
	"layeh.com/radius/rfc2868"
	"layeh.com/radius/rfc2869"
	"layeh.com/radius/rfc3162"
	"layeh.com/radius/rfc3576"
	"layeh.com/radius/rfc4072"
	"layeh.com/radius/rfc4372"
	"layeh.com/radius/rfc4675"
	"layeh.com/radius/rfc4818"
	"layeh.com/radius/rfc4849"
	"layeh.com/radius/rfc5090"
	"layeh.com/radius/rfc5447"
	"layeh.com/radius/rfc5580"
	"layeh.com/radius/rfc5607"
	"layeh.com/radius/rfc5904"
	"layeh.com/radius/rfc6519"
	"layeh.com/radius/rfc6572"
	"layeh.com/radius/rfc6677"
	"layeh.com/radius/rfc6911"
	"layeh.com/radius/rfc7055"
	"layeh.com/radius/rfc7268"
)

var RadiusTypeName = map[radius.Type]string{
	rfc5904.PKMSSCert_Type:                      "PKMSSCert",
	rfc5904.PKMCACert_Type:                      "PKMCACert",
	rfc5904.PKMConfigSettings_Type:              "PKMConfigSettings",
	rfc5904.PKMCryptosuiteList_Type:             "PKMCryptosuiteList",
	rfc5904.PKMSAID_Type:                        "PKMSAID",
	rfc5904.PKMSADescriptor_Type:                "PKMSADescriptor",
	rfc5904.PKMAuthKey_Type:                     "PKMAuthKey",
	rfc6911.FramedIPv6Address_Type:              "FramedIPv6Address",
	rfc6911.DNSServerIPv6Address_Type:           "DNSServerIPv6Address",
	rfc6911.RouteIPv6Information_Type:           "RouteIPv6Information",
	rfc6911.DelegatedIPv6PrefixPool_Type:        "DelegatedIPv6PrefixPool",
	rfc6911.StatefulIPv6AddressPool_Type:        "StatefulIPv6AddressPool",
	rfc5447.MIP6FeatureVector_Type:              "MIP6FeatureVector",
	rfc5447.MIP6HomeLinkPrefix_Type:             "MIP6HomeLinkPrefix",
	rfc5580.OperatorName_Type:                   "OperatorName",
	rfc5580.LocationInformation_Type:            "LocationInformation",
	rfc5580.LocationData_Type:                   "LocationData",
	rfc5580.BasicLocationPolicyRules_Type:       "BasicLocationPolicyRules",
	rfc5580.ExtendedLocationPolicyRules_Type:    "ExtendedLocationPolicyRules",
	rfc5580.LocationCapable_Type:                "LocationCapable",
	rfc5580.RequestedLocationInfo_Type:          "RequestedLocationInfo",
	rfc2869.AcctInputGigawords_Type:             "AcctInputGigawords",
	rfc2869.AcctOutputGigawords_Type:            "AcctOutputGigawords",
	rfc2869.EventTimestamp_Type:                 "EventTimestamp",
	rfc2869.ARAPPassword_Type:                   "ARAPPassword",
	rfc2869.ARAPFeatures_Type:                   "ARAPFeatures",
	rfc2869.ARAPZoneAccess_Type:                 "ARAPZoneAccess",
	rfc2869.ARAPSecurity_Type:                   "ARAPSecurity",
	rfc2869.ARAPSecurityData_Type:               "ARAPSecurityData",
	rfc2869.PasswordRetry_Type:                  "PasswordRetry",
	rfc2869.Prompt_Type:                         "Prompt",
	rfc2869.ConnectInfo_Type:                    "ConnectInfo",
	rfc2869.ConfigurationToken_Type:             "ConfigurationToken",
	rfc2869.EAPMessage_Type:                     "EAPMessage",
	rfc2869.MessageAuthenticator_Type:           "MessageAuthenticator",
	rfc2869.ARAPChallengeResponse_Type:          "ARAPChallengeResponse",
	rfc2869.AcctInterimInterval_Type:            "AcctInterimInterval",
	rfc2869.NASPortID_Type:                      "NASPortID",
	rfc2869.FramedPool_Type:                     "FramedPool",
	rfc4072.EAPKeyName_Type:                     "EAPKeyName",
	rfc4372.ChargeableUserIdentity_Type:         "ChargeableUserIdentity",
	rfc4849.NASFilterRule_Type:                  "NASFilterRule",
	rfc7268.AllowedCalledStationID_Type:         "AllowedCalledStationID",
	rfc7268.EAPPeerID_Type:                      "EAPPeerID",
	rfc7268.EAPServerID_Type:                    "EAPServerID",
	rfc7268.MobilityDomainID_Type:               "MobilityDomainID",
	rfc7268.PreauthTimeout_Type:                 "PreauthTimeout",
	rfc7268.NetworkIDName_Type:                  "NetworkIDName",
	rfc7268.EAPoLAnnouncement_Type:              "EAPoLAnnouncement",
	rfc7268.WLANHESSID_Type:                     "WLANHESSID",
	rfc7268.WLANVenueInfo_Type:                  "WLANVenueInfo",
	rfc7268.WLANVenueLanguage_Type:              "WLANVenueLanguage",
	rfc7268.WLANVenueName_Type:                  "WLANVenueName",
	rfc7268.WLANReasonCode_Type:                 "WLANReasonCode",
	rfc7268.WLANPairwiseCipher_Type:             "WLANPairwiseCipher",
	rfc7268.WLANGroupCipher_Type:                "WLANGroupCipher",
	rfc7268.WLANAKMSuite_Type:                   "WLANAKMSuite",
	rfc7268.WLANGroupMgmtCipher_Type:            "WLANGroupMgmtCipher",
	rfc7268.WLANRFBand_Type:                     "WLANRFBand",
	rfc3162.NASIPv6Address_Type:                 "NASIPv6Address",
	rfc3162.FramedInterfaceID_Type:              "FramedInterfaceID",
	rfc3162.FramedIPv6Prefix_Type:               "FramedIPv6Prefix",
	rfc3162.LoginIPv6Host_Type:                  "LoginIPv6Host",
	rfc3162.FramedIPv6Route_Type:                "FramedIPv6Route",
	rfc3162.FramedIPv6Pool_Type:                 "FramedIPv6Pool",
	rfc2865.UserName_Type:                       "UserName",
	rfc2865.UserPassword_Type:                   "UserPassword",
	rfc2865.CHAPPassword_Type:                   "CHAPPassword",
	rfc2865.NASIPAddress_Type:                   "NASIPAddress",
	rfc2865.NASPort_Type:                        "NASPort",
	rfc2865.ServiceType_Type:                    "ServiceType",
	rfc2865.FramedProtocol_Type:                 "FramedProtocol",
	rfc2865.FramedIPAddress_Type:                "FramedIPAddress",
	rfc2865.FramedIPNetmask_Type:                "FramedIPNetmask",
	rfc2865.FramedRouting_Type:                  "FramedRouting",
	rfc2865.FilterID_Type:                       "FilterID",
	rfc2865.FramedMTU_Type:                      "FramedMTU",
	rfc2865.FramedCompression_Type:              "FramedCompression",
	rfc2865.LoginIPHost_Type:                    "LoginIPHost",
	rfc2865.LoginService_Type:                   "LoginService",
	rfc2865.LoginTCPPort_Type:                   "LoginTCPPort",
	rfc2865.ReplyMessage_Type:                   "ReplyMessage",
	rfc2865.CallbackNumber_Type:                 "CallbackNumber",
	rfc2865.CallbackID_Type:                     "CallbackID",
	rfc2865.FramedRoute_Type:                    "FramedRoute",
	rfc2865.FramedIPXNetwork_Type:               "FramedIPXNetwork",
	rfc2865.State_Type:                          "State",
	rfc2865.Class_Type:                          "Class",
	rfc2865.VendorSpecific_Type:                 "VendorSpecific",
	rfc2865.SessionTimeout_Type:                 "SessionTimeout",
	rfc2865.IdleTimeout_Type:                    "IdleTimeout",
	rfc2865.TerminationAction_Type:              "TerminationAction",
	rfc2865.CalledStationID_Type:                "CalledStationID",
	rfc2865.CallingStationID_Type:               "CallingStationID",
	rfc2865.NASIdentifier_Type:                  "NASIdentifier",
	rfc2865.ProxyState_Type:                     "ProxyState",
	rfc2865.LoginLATService_Type:                "LoginLATService",
	rfc2865.LoginLATNode_Type:                   "LoginLATNode",
	rfc2865.LoginLATGroup_Type:                  "LoginLATGroup",
	rfc2865.FramedAppleTalkLink_Type:            "FramedAppleTalkLink",
	rfc2865.FramedAppleTalkNetwork_Type:         "FramedAppleTalkNetwork",
	rfc2865.FramedAppleTalkZone_Type:            "FramedAppleTalkZone",
	rfc2865.CHAPChallenge_Type:                  "CHAPChallenge",
	rfc2865.NASPortType_Type:                    "NASPortType",
	rfc2865.PortLimit_Type:                      "PortLimit",
	rfc2865.LoginLATPort_Type:                   "LoginLATPort",
	rfc6572.MobileNodeIdentifier_Type:           "MobileNodeIdentifier",
	rfc6572.ServiceSelection_Type:               "ServiceSelection",
	rfc6572.PMIP6HomeLMAIPv6Address_Type:        "PMIP6HomeLMAIPv6Address",
	rfc6572.PMIP6VisitedLMAIPv6Address_Type:     "PMIP6VisitedLMAIPv6Address",
	rfc6572.PMIP6HomeLMAIPv4Address_Type:        "PMIP6HomeLMAIPv4Address",
	rfc6572.PMIP6VisitedLMAIPv4Address_Type:     "PMIP6VisitedLMAIPv4Address",
	rfc6572.PMIP6HomeHNPrefix_Type:              "PMIP6HomeHNPrefix",
	rfc6572.PMIP6VisitedHNPrefix_Type:           "PMIP6VisitedHNPrefix",
	rfc6572.PMIP6HomeInterfaceID_Type:           "PMIP6HomeInterfaceID",
	rfc6572.PMIP6VisitedInterfaceID_Type:        "PMIP6VisitedInterfaceID",
	rfc6572.PMIP6HomeDHCP4ServerAddress_Type:    "PMIP6HomeDHCP4ServerAddress",
	rfc6572.PMIP6VisitedDHCP4ServerAddress_Type: "PMIP6VisitedDHCP4ServerAddress",
	rfc6572.PMIP6HomeDHCP6ServerAddress_Type:    "PMIP6HomeDHCP6ServerAddress",
	rfc6572.PMIP6VisitedDHCP6ServerAddress_Type: "PMIP6VisitedDHCP6ServerAddress",
	rfc6572.PMIP6HomeIPv4Gateway_Type:           "PMIP6HomeIPv4Gateway",
	rfc6572.PMIP6VisitedIPv4Gateway_Type:        "PMIP6VisitedIPv4Gateway",
	rfc4818.DelegatedIPv6Prefix_Type:            "DelegatedIPv6Prefix",
	rfc6519.DSLiteTunnelName_Type:               "DSLiteTunnelName",
	rfc7055.GSSAcceptorServiceName_Type:         "GSSAcceptorServiceName",
	rfc7055.GSSAcceptorHostName_Type:            "GSSAcceptorHostName",
	rfc7055.GSSAcceptorServiceSpecifics_Type:    "GSSAcceptorServiceSpecifics",
	rfc7055.GSSAcceptorRealmName_Type:           "GSSAcceptorRealmName",
	rfc2867.AcctTunnelConnection_Type:           "AcctTunnelConnection",
	rfc2867.AcctTunnelPacketsLost_Type:          "AcctTunnelPacketsLost",
	rfc3576.ErrorCause_Type:                     "ErrorCause",
	rfc6677.EAPLowerLayer_Type:                  "EAPLowerLayer",
	rfc2866.AcctStatusType_Type:                 "AcctStatusType",
	rfc2866.AcctDelayTime_Type:                  "AcctDelayTime",
	rfc2866.AcctInputOctets_Type:                "AcctInputOctets",
	rfc2866.AcctOutputOctets_Type:               "AcctOutputOctets",
	rfc2866.AcctSessionID_Type:                  "AcctSessionID",
	rfc2866.AcctAuthentic_Type:                  "AcctAuthentic",
	rfc2866.AcctSessionTime_Type:                "AcctSessionTime",
	rfc2866.AcctInputPackets_Type:               "AcctInputPackets",
	rfc2866.AcctOutputPackets_Type:              "AcctOutputPackets",
	rfc2866.AcctTerminateCause_Type:             "AcctTerminateCause",
	rfc2866.AcctMultiSessionID_Type:             "AcctMultiSessionID",
	rfc2866.AcctLinkCount_Type:                  "AcctLinkCount",
	rfc4675.EgressVLANID_Type:                   "EgressVLANID",
	rfc4675.IngressFilters_Type:                 "IngressFilters",
	rfc4675.EgressVLANName_Type:                 "EgressVLANName",
	rfc4675.UserPriorityTable_Type:              "UserPriorityTable",
	rfc5607.FramedManagement_Type:               "FramedManagement",
	rfc5607.ManagementTransportProtection_Type:  "ManagementTransportProtection",
	rfc5607.ManagementPolicyID_Type:             "ManagementPolicyID",
	rfc5607.ManagementPrivilegeLevel_Type:       "ManagementPrivilegeLevel",
	rfc5090.DigestResponse_Type:                 "DigestResponse",
	rfc5090.DigestRealm_Type:                    "DigestRealm",
	rfc5090.DigestNonce_Type:                    "DigestNonce",
	rfc5090.DigestResponseAuth_Type:             "DigestResponseAuth",
	rfc5090.DigestNextnonce_Type:                "DigestNextnonce",
	rfc5090.DigestMethod_Type:                   "DigestMethod",
	rfc5090.DigestURI_Type:                      "DigestURI",
	rfc5090.DigestQop_Type:                      "DigestQop",
	rfc5090.DigestAlgorithm_Type:                "DigestAlgorithm",
	rfc5090.DigestEntityBodyHash_Type:           "DigestEntityBodyHash",
	rfc5090.DigestCNonce_Type:                   "DigestCNonce",
	rfc5090.DigestNonceCount_Type:               "DigestNonceCount",
	rfc5090.DigestUsername_Type:                 "DigestUsername",
	rfc5090.DigestOpaque_Type:                   "DigestOpaque",
	rfc5090.DigestAuthParam_Type:                "DigestAuthParam",
	rfc5090.DigestAKAAuts_Type:                  "DigestAKAAuts",
	rfc5090.DigestDomain_Type:                   "DigestDomain",
	rfc5090.DigestStale_Type:                    "DigestStale",
	rfc5090.DigestHA1_Type:                      "DigestHA1",
	rfc5090.SIPAOR_Type:                         "SIPAOR",
	rfc2868.TunnelType_Type:                     "TunnelType",
	rfc2868.TunnelMediumType_Type:               "TunnelMediumType",
	rfc2868.TunnelClientEndpoint_Type:           "TunnelClientEndpoint",
	rfc2868.TunnelServerEndpoint_Type:           "TunnelServerEndpoint",
	rfc2868.TunnelPassword_Type:                 "TunnelPassword",
	rfc2868.TunnelPrivateGroupID_Type:           "TunnelPrivateGroupID",
	rfc2868.TunnelAssignmentID_Type:             "TunnelAssignmentID",
	rfc2868.TunnelPreference_Type:               "TunnelPreference",
	rfc2868.TunnelClientAuthID_Type:             "TunnelClientAuthID",
	rfc2868.TunnelServerAuthID_Type:             "TunnelServerAuthID",
}
