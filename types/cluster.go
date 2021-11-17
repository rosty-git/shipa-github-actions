package types

import (
    "encoding/json"
    "fmt"
    "github.com/brunoa19/shipa-github-actions/shipa"
    "io/ioutil"
    "os"
)

type Cluster struct {
    Name      string            `json:"name" yaml:"name"`
    Endpoint  *ClusterEndpoint  `json:"endpoint" yaml:"endpoint"`
    Resources *ClusterResources `json:"resources,omitempty" yaml:"resources,omitempty"`
}

func (c *Cluster) ToShipaCluster() (*shipa.Cluster, error) {
    var frameworks []*shipa.Framework
    if c.Resources != nil && c.Resources.Frameworks != nil {
        for _, name := range c.Resources.Frameworks.Name {
            frameworks = append(frameworks, &shipa.Framework{
                Name: name,
            })
        }
        c.Resources.Frameworks = nil
    }

    rawJson, err := json.Marshal(c)
    if err != nil {
        return nil, err
    }

    cluster := &shipa.Cluster{}
    err = json.Unmarshal(rawJson, cluster)
    if err != nil {
        return nil, err
    }

    if frameworks != nil {
        if cluster.Resources == nil {
            cluster.Resources = &shipa.ClusterResources{}
        }
        cluster.Resources.Frameworks = frameworks
    }

    if cluster.Endpoint != nil {
        cluster.Endpoint.Token = useFileOrValue(cluster.Endpoint.Token)
        cluster.Endpoint.Certificate = useFileOrValue(cluster.Endpoint.Certificate)
        cluster.Endpoint.ClientCertificate = useFileOrValue(cluster.Endpoint.ClientCertificate)
        cluster.Endpoint.ClientKey = useFileOrValue(cluster.Endpoint.ClientKey)
    }

    return cluster, nil
}

func useFileOrValue(value string) string {
    if value == "" {
        return ""
    }

    data, err := readFile(value)
    if err != nil {
        return value
    }
    return string(data)
}

func readFile(path string) ([]byte, error) {
    if _, err := os.Stat(path); err != nil {
        return nil, fmt.Errorf("invalid file path: %v", err)
    }

    bytes, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %v", err)
    }

    return bytes, nil
}

// ClusterEndpoint - part of Cluster object
type ClusterEndpoint struct {
    Addresses         []string `json:"addresses,omitempty" yaml:"addresses,omitempty"`
    Certificate       string   `json:"caCert,omitempty" yaml:"caCert,omitempty"`
    ClientCertificate string   `json:"clientCert,omitempty" yaml:"clientCert,omitempty"`
    ClientKey         string   `json:"clientKey,omitempty" yaml:"clientKey,omitempty"`
    Token             string   `json:"token,omitempty" yaml:"token,omitempty"`
    Username          string   `json:"username,omitempty" yaml:"username,omitempty"`
    Password          string   `json:"password,omitempty" yaml:"password,omitempty"`
}

// ClusterResources - part of Cluster object
type ClusterResources struct {
    Frameworks         *Framework         `json:"frameworks,omitempty" yaml:"frameworks,omitempty"`
    IngressControllers []*IngressController `json:"ingressControllers,omitempty" yaml:"ingressControllers,omitempty"`
}

// IngressController - part of ClusterResources object
type IngressController struct {
    IngressIP     string `json:"ingressIp,omitempty" yaml:"ingressIp,omitempty"`
    ServiceType   string `json:"serviceType,omitempty" yaml:"serviceType,omitempty"`
    Type          string `json:"type,omitempty" yaml:"type,omitempty"`
    HTTPPort      int64  `json:"httpPort,omitempty" yaml:"httpPort,omitempty"`
    HTTPSPort     int64  `json:"httpsPort,omitempty" yaml:"httpsPort,omitempty"`
    ProtectedPort int64  `json:"protectedPort,omitempty" yaml:"protectedPort,omitempty"`
    Debug         bool   `json:"debug" yaml:"debug"`
    AcmeEmail     string `json:"acmeEmail,omitempty" yaml:"acmeEmail,omitempty"`
    AcmeServer    string `json:"acmeServer,omitempty" yaml:"acmeServer,omitempty"`
}

// Framework - part of ClusterResources object
type Framework struct {
    Name []string `json:"name,omitempty" yaml:"name,omitempty"`
}
