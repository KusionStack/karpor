//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Karpor Authors.

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1beta1

import (
	url "net/url"
	unsafe "unsafe"

	cluster "github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Cluster)(nil), (*cluster.Cluster)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Cluster_To_cluster_Cluster(a.(*Cluster), b.(*cluster.Cluster), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.Cluster)(nil), (*Cluster)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_Cluster_To_v1beta1_Cluster(a.(*cluster.Cluster), b.(*Cluster), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterAccess)(nil), (*cluster.ClusterAccess)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterAccess_To_cluster_ClusterAccess(a.(*ClusterAccess), b.(*cluster.ClusterAccess), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterAccess)(nil), (*ClusterAccess)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterAccess_To_v1beta1_ClusterAccess(a.(*cluster.ClusterAccess), b.(*ClusterAccess), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterAccessCredential)(nil), (*cluster.ClusterAccessCredential)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterAccessCredential_To_cluster_ClusterAccessCredential(a.(*ClusterAccessCredential), b.(*cluster.ClusterAccessCredential), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterAccessCredential)(nil), (*ClusterAccessCredential)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterAccessCredential_To_v1beta1_ClusterAccessCredential(a.(*cluster.ClusterAccessCredential), b.(*ClusterAccessCredential), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterList)(nil), (*cluster.ClusterList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterList_To_cluster_ClusterList(a.(*ClusterList), b.(*cluster.ClusterList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterList)(nil), (*ClusterList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterList_To_v1beta1_ClusterList(a.(*cluster.ClusterList), b.(*ClusterList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterProxyOptions)(nil), (*cluster.ClusterProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterProxyOptions_To_cluster_ClusterProxyOptions(a.(*ClusterProxyOptions), b.(*cluster.ClusterProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterProxyOptions)(nil), (*ClusterProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterProxyOptions_To_v1beta1_ClusterProxyOptions(a.(*cluster.ClusterProxyOptions), b.(*ClusterProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterSpec)(nil), (*cluster.ClusterSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterSpec_To_cluster_ClusterSpec(a.(*ClusterSpec), b.(*cluster.ClusterSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterSpec)(nil), (*ClusterSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterSpec_To_v1beta1_ClusterSpec(a.(*cluster.ClusterSpec), b.(*ClusterSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterStatus)(nil), (*cluster.ClusterStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterStatus_To_cluster_ClusterStatus(a.(*ClusterStatus), b.(*cluster.ClusterStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterStatus)(nil), (*ClusterStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterStatus_To_v1beta1_ClusterStatus(a.(*cluster.ClusterStatus), b.(*ClusterStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ExecConfig)(nil), (*cluster.ExecConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ExecConfig_To_cluster_ExecConfig(a.(*ExecConfig), b.(*cluster.ExecConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ExecConfig)(nil), (*ExecConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ExecConfig_To_v1beta1_ExecConfig(a.(*cluster.ExecConfig), b.(*ExecConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ExecEnvVar)(nil), (*cluster.ExecEnvVar)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ExecEnvVar_To_cluster_ExecEnvVar(a.(*ExecEnvVar), b.(*cluster.ExecEnvVar), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ExecEnvVar)(nil), (*ExecEnvVar)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ExecEnvVar_To_v1beta1_ExecEnvVar(a.(*cluster.ExecEnvVar), b.(*ExecEnvVar), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*X509)(nil), (*cluster.X509)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_X509_To_cluster_X509(a.(*X509), b.(*cluster.X509), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.X509)(nil), (*X509)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_X509_To_v1beta1_X509(a.(*cluster.X509), b.(*X509), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*url.Values)(nil), (*ClusterProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_url_Values_To_v1beta1_ClusterProxyOptions(a.(*url.Values), b.(*ClusterProxyOptions), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1beta1_Cluster_To_cluster_Cluster(in *Cluster, out *cluster.Cluster, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta1_ClusterSpec_To_cluster_ClusterSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_ClusterStatus_To_cluster_ClusterStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_Cluster_To_cluster_Cluster is an autogenerated conversion function.
func Convert_v1beta1_Cluster_To_cluster_Cluster(in *Cluster, out *cluster.Cluster, s conversion.Scope) error {
	return autoConvert_v1beta1_Cluster_To_cluster_Cluster(in, out, s)
}

func autoConvert_cluster_Cluster_To_v1beta1_Cluster(in *cluster.Cluster, out *Cluster, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_cluster_ClusterSpec_To_v1beta1_ClusterSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_cluster_ClusterStatus_To_v1beta1_ClusterStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_cluster_Cluster_To_v1beta1_Cluster is an autogenerated conversion function.
func Convert_cluster_Cluster_To_v1beta1_Cluster(in *cluster.Cluster, out *Cluster, s conversion.Scope) error {
	return autoConvert_cluster_Cluster_To_v1beta1_Cluster(in, out, s)
}

func autoConvert_v1beta1_ClusterAccess_To_cluster_ClusterAccess(in *ClusterAccess, out *cluster.ClusterAccess, s conversion.Scope) error {
	out.Endpoint = in.Endpoint
	out.CABundle = *(*[]byte)(unsafe.Pointer(&in.CABundle))
	out.Insecure = (*bool)(unsafe.Pointer(in.Insecure))
	out.Credential = (*cluster.ClusterAccessCredential)(unsafe.Pointer(in.Credential))
	return nil
}

// Convert_v1beta1_ClusterAccess_To_cluster_ClusterAccess is an autogenerated conversion function.
func Convert_v1beta1_ClusterAccess_To_cluster_ClusterAccess(in *ClusterAccess, out *cluster.ClusterAccess, s conversion.Scope) error {
	return autoConvert_v1beta1_ClusterAccess_To_cluster_ClusterAccess(in, out, s)
}

func autoConvert_cluster_ClusterAccess_To_v1beta1_ClusterAccess(in *cluster.ClusterAccess, out *ClusterAccess, s conversion.Scope) error {
	out.Endpoint = in.Endpoint
	out.CABundle = *(*[]byte)(unsafe.Pointer(&in.CABundle))
	out.Insecure = (*bool)(unsafe.Pointer(in.Insecure))
	out.Credential = (*ClusterAccessCredential)(unsafe.Pointer(in.Credential))
	return nil
}

// Convert_cluster_ClusterAccess_To_v1beta1_ClusterAccess is an autogenerated conversion function.
func Convert_cluster_ClusterAccess_To_v1beta1_ClusterAccess(in *cluster.ClusterAccess, out *ClusterAccess, s conversion.Scope) error {
	return autoConvert_cluster_ClusterAccess_To_v1beta1_ClusterAccess(in, out, s)
}

func autoConvert_v1beta1_ClusterAccessCredential_To_cluster_ClusterAccessCredential(in *ClusterAccessCredential, out *cluster.ClusterAccessCredential, s conversion.Scope) error {
	out.Type = cluster.CredentialType(in.Type)
	out.ServiceAccountToken = in.ServiceAccountToken
	out.X509 = (*cluster.X509)(unsafe.Pointer(in.X509))
	out.ExecConfig = (*cluster.ExecConfig)(unsafe.Pointer(in.ExecConfig))
	return nil
}

// Convert_v1beta1_ClusterAccessCredential_To_cluster_ClusterAccessCredential is an autogenerated conversion function.
func Convert_v1beta1_ClusterAccessCredential_To_cluster_ClusterAccessCredential(in *ClusterAccessCredential, out *cluster.ClusterAccessCredential, s conversion.Scope) error {
	return autoConvert_v1beta1_ClusterAccessCredential_To_cluster_ClusterAccessCredential(in, out, s)
}

func autoConvert_cluster_ClusterAccessCredential_To_v1beta1_ClusterAccessCredential(in *cluster.ClusterAccessCredential, out *ClusterAccessCredential, s conversion.Scope) error {
	out.Type = CredentialType(in.Type)
	out.ServiceAccountToken = in.ServiceAccountToken
	out.X509 = (*X509)(unsafe.Pointer(in.X509))
	out.ExecConfig = (*ExecConfig)(unsafe.Pointer(in.ExecConfig))
	return nil
}

// Convert_cluster_ClusterAccessCredential_To_v1beta1_ClusterAccessCredential is an autogenerated conversion function.
func Convert_cluster_ClusterAccessCredential_To_v1beta1_ClusterAccessCredential(in *cluster.ClusterAccessCredential, out *ClusterAccessCredential, s conversion.Scope) error {
	return autoConvert_cluster_ClusterAccessCredential_To_v1beta1_ClusterAccessCredential(in, out, s)
}

func autoConvert_v1beta1_ClusterList_To_cluster_ClusterList(in *ClusterList, out *cluster.ClusterList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]cluster.Cluster)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1beta1_ClusterList_To_cluster_ClusterList is an autogenerated conversion function.
func Convert_v1beta1_ClusterList_To_cluster_ClusterList(in *ClusterList, out *cluster.ClusterList, s conversion.Scope) error {
	return autoConvert_v1beta1_ClusterList_To_cluster_ClusterList(in, out, s)
}

func autoConvert_cluster_ClusterList_To_v1beta1_ClusterList(in *cluster.ClusterList, out *ClusterList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]Cluster)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_cluster_ClusterList_To_v1beta1_ClusterList is an autogenerated conversion function.
func Convert_cluster_ClusterList_To_v1beta1_ClusterList(in *cluster.ClusterList, out *ClusterList, s conversion.Scope) error {
	return autoConvert_cluster_ClusterList_To_v1beta1_ClusterList(in, out, s)
}

func autoConvert_v1beta1_ClusterProxyOptions_To_cluster_ClusterProxyOptions(in *ClusterProxyOptions, out *cluster.ClusterProxyOptions, s conversion.Scope) error {
	out.Path = in.Path
	return nil
}

// Convert_v1beta1_ClusterProxyOptions_To_cluster_ClusterProxyOptions is an autogenerated conversion function.
func Convert_v1beta1_ClusterProxyOptions_To_cluster_ClusterProxyOptions(in *ClusterProxyOptions, out *cluster.ClusterProxyOptions, s conversion.Scope) error {
	return autoConvert_v1beta1_ClusterProxyOptions_To_cluster_ClusterProxyOptions(in, out, s)
}

func autoConvert_cluster_ClusterProxyOptions_To_v1beta1_ClusterProxyOptions(in *cluster.ClusterProxyOptions, out *ClusterProxyOptions, s conversion.Scope) error {
	out.Path = in.Path
	return nil
}

// Convert_cluster_ClusterProxyOptions_To_v1beta1_ClusterProxyOptions is an autogenerated conversion function.
func Convert_cluster_ClusterProxyOptions_To_v1beta1_ClusterProxyOptions(in *cluster.ClusterProxyOptions, out *ClusterProxyOptions, s conversion.Scope) error {
	return autoConvert_cluster_ClusterProxyOptions_To_v1beta1_ClusterProxyOptions(in, out, s)
}

func autoConvert_url_Values_To_v1beta1_ClusterProxyOptions(in *url.Values, out *ClusterProxyOptions, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	if values, ok := map[string][]string(*in)["path"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_string(&values, &out.Path, s); err != nil {
			return err
		}
	} else {
		out.Path = ""
	}
	return nil
}

// Convert_url_Values_To_v1beta1_ClusterProxyOptions is an autogenerated conversion function.
func Convert_url_Values_To_v1beta1_ClusterProxyOptions(in *url.Values, out *ClusterProxyOptions, s conversion.Scope) error {
	return autoConvert_url_Values_To_v1beta1_ClusterProxyOptions(in, out, s)
}

func autoConvert_v1beta1_ClusterSpec_To_cluster_ClusterSpec(in *ClusterSpec, out *cluster.ClusterSpec, s conversion.Scope) error {
	out.Provider = in.Provider
	if err := Convert_v1beta1_ClusterAccess_To_cluster_ClusterAccess(&in.Access, &out.Access, s); err != nil {
		return err
	}
	out.Description = in.Description
	out.DisplayName = in.DisplayName
	out.Finalized = (*bool)(unsafe.Pointer(in.Finalized))
	return nil
}

// Convert_v1beta1_ClusterSpec_To_cluster_ClusterSpec is an autogenerated conversion function.
func Convert_v1beta1_ClusterSpec_To_cluster_ClusterSpec(in *ClusterSpec, out *cluster.ClusterSpec, s conversion.Scope) error {
	return autoConvert_v1beta1_ClusterSpec_To_cluster_ClusterSpec(in, out, s)
}

func autoConvert_cluster_ClusterSpec_To_v1beta1_ClusterSpec(in *cluster.ClusterSpec, out *ClusterSpec, s conversion.Scope) error {
	out.Provider = in.Provider
	if err := Convert_cluster_ClusterAccess_To_v1beta1_ClusterAccess(&in.Access, &out.Access, s); err != nil {
		return err
	}
	out.Description = in.Description
	out.DisplayName = in.DisplayName
	out.Finalized = (*bool)(unsafe.Pointer(in.Finalized))
	return nil
}

// Convert_cluster_ClusterSpec_To_v1beta1_ClusterSpec is an autogenerated conversion function.
func Convert_cluster_ClusterSpec_To_v1beta1_ClusterSpec(in *cluster.ClusterSpec, out *ClusterSpec, s conversion.Scope) error {
	return autoConvert_cluster_ClusterSpec_To_v1beta1_ClusterSpec(in, out, s)
}

func autoConvert_v1beta1_ClusterStatus_To_cluster_ClusterStatus(in *ClusterStatus, out *cluster.ClusterStatus, s conversion.Scope) error {
	out.Healthy = in.Healthy
	return nil
}

// Convert_v1beta1_ClusterStatus_To_cluster_ClusterStatus is an autogenerated conversion function.
func Convert_v1beta1_ClusterStatus_To_cluster_ClusterStatus(in *ClusterStatus, out *cluster.ClusterStatus, s conversion.Scope) error {
	return autoConvert_v1beta1_ClusterStatus_To_cluster_ClusterStatus(in, out, s)
}

func autoConvert_cluster_ClusterStatus_To_v1beta1_ClusterStatus(in *cluster.ClusterStatus, out *ClusterStatus, s conversion.Scope) error {
	out.Healthy = in.Healthy
	return nil
}

// Convert_cluster_ClusterStatus_To_v1beta1_ClusterStatus is an autogenerated conversion function.
func Convert_cluster_ClusterStatus_To_v1beta1_ClusterStatus(in *cluster.ClusterStatus, out *ClusterStatus, s conversion.Scope) error {
	return autoConvert_cluster_ClusterStatus_To_v1beta1_ClusterStatus(in, out, s)
}

func autoConvert_v1beta1_ExecConfig_To_cluster_ExecConfig(in *ExecConfig, out *cluster.ExecConfig, s conversion.Scope) error {
	out.Command = in.Command
	out.Args = *(*[]string)(unsafe.Pointer(&in.Args))
	out.Env = *(*[]cluster.ExecEnvVar)(unsafe.Pointer(&in.Env))
	out.APIVersion = in.APIVersion
	out.InstallHint = in.InstallHint
	out.ProvideClusterInfo = in.ProvideClusterInfo
	out.InteractiveMode = in.InteractiveMode
	return nil
}

// Convert_v1beta1_ExecConfig_To_cluster_ExecConfig is an autogenerated conversion function.
func Convert_v1beta1_ExecConfig_To_cluster_ExecConfig(in *ExecConfig, out *cluster.ExecConfig, s conversion.Scope) error {
	return autoConvert_v1beta1_ExecConfig_To_cluster_ExecConfig(in, out, s)
}

func autoConvert_cluster_ExecConfig_To_v1beta1_ExecConfig(in *cluster.ExecConfig, out *ExecConfig, s conversion.Scope) error {
	out.Command = in.Command
	out.Args = *(*[]string)(unsafe.Pointer(&in.Args))
	out.Env = *(*[]ExecEnvVar)(unsafe.Pointer(&in.Env))
	out.APIVersion = in.APIVersion
	out.InstallHint = in.InstallHint
	out.ProvideClusterInfo = in.ProvideClusterInfo
	out.InteractiveMode = in.InteractiveMode
	return nil
}

// Convert_cluster_ExecConfig_To_v1beta1_ExecConfig is an autogenerated conversion function.
func Convert_cluster_ExecConfig_To_v1beta1_ExecConfig(in *cluster.ExecConfig, out *ExecConfig, s conversion.Scope) error {
	return autoConvert_cluster_ExecConfig_To_v1beta1_ExecConfig(in, out, s)
}

func autoConvert_v1beta1_ExecEnvVar_To_cluster_ExecEnvVar(in *ExecEnvVar, out *cluster.ExecEnvVar, s conversion.Scope) error {
	out.Name = in.Name
	out.Value = in.Value
	return nil
}

// Convert_v1beta1_ExecEnvVar_To_cluster_ExecEnvVar is an autogenerated conversion function.
func Convert_v1beta1_ExecEnvVar_To_cluster_ExecEnvVar(in *ExecEnvVar, out *cluster.ExecEnvVar, s conversion.Scope) error {
	return autoConvert_v1beta1_ExecEnvVar_To_cluster_ExecEnvVar(in, out, s)
}

func autoConvert_cluster_ExecEnvVar_To_v1beta1_ExecEnvVar(in *cluster.ExecEnvVar, out *ExecEnvVar, s conversion.Scope) error {
	out.Name = in.Name
	out.Value = in.Value
	return nil
}

// Convert_cluster_ExecEnvVar_To_v1beta1_ExecEnvVar is an autogenerated conversion function.
func Convert_cluster_ExecEnvVar_To_v1beta1_ExecEnvVar(in *cluster.ExecEnvVar, out *ExecEnvVar, s conversion.Scope) error {
	return autoConvert_cluster_ExecEnvVar_To_v1beta1_ExecEnvVar(in, out, s)
}

func autoConvert_v1beta1_X509_To_cluster_X509(in *X509, out *cluster.X509, s conversion.Scope) error {
	out.Certificate = *(*[]byte)(unsafe.Pointer(&in.Certificate))
	out.PrivateKey = *(*[]byte)(unsafe.Pointer(&in.PrivateKey))
	return nil
}

// Convert_v1beta1_X509_To_cluster_X509 is an autogenerated conversion function.
func Convert_v1beta1_X509_To_cluster_X509(in *X509, out *cluster.X509, s conversion.Scope) error {
	return autoConvert_v1beta1_X509_To_cluster_X509(in, out, s)
}

func autoConvert_cluster_X509_To_v1beta1_X509(in *cluster.X509, out *X509, s conversion.Scope) error {
	out.Certificate = *(*[]byte)(unsafe.Pointer(&in.Certificate))
	out.PrivateKey = *(*[]byte)(unsafe.Pointer(&in.PrivateKey))
	return nil
}

// Convert_cluster_X509_To_v1beta1_X509 is an autogenerated conversion function.
func Convert_cluster_X509_To_v1beta1_X509(in *cluster.X509, out *X509, s conversion.Scope) error {
	return autoConvert_cluster_X509_To_v1beta1_X509(in, out, s)
}
