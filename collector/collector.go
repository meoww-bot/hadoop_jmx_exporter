package collector

import (
	"encoding/json"
	"fmt"
	"hadoop_jmx_exporter/lib"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/prometheus/client_golang/prometheus"
)

type CollectorFunc func(target Target, registry *prometheus.Registry) bool

type Target struct {
	Url           string
	ExporterName  string
	BodyData      []byte
	KrbAuthMethod string
	KrbPrincipal  string
	KrbPassword   string
	KrbKtPath     string
	Logger        log.Logger
}

var (
	Collectors = map[string]CollectorFunc{
		"NameNode":          NameNodeCollector,
		"DataNode":          DataNodeCollector,
		"ResourceManager":   ResourceManagerCollector,
		"JournalNode":       JournalNodeCollector,
		"hiveserver2":       HiveServer2Collector,
		"NodeManager":       NodeManagerCollector,
		"HbaseMaster":       HbaseMasterCollector,
		"HbaseRegionServer": HbaseRegionServerCollector,
	}
)

func (t *Target) getCollectorName() error {

	var data []byte
	var err error

	UrlHostname, err := lib.ExtractDomainFromURL(t.Url)

	if err != nil {
		level.Error(t.Logger).Log("msg", "Error extract domain from url", "err", err)
		return err
	}

	if UrlHostname == "127.0.0.1" {
		resp, err := http.Get(t.Url)
		if err != nil {
			level.Error(t.Logger).Log("msg", "Error get url", "err", err)

			return err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			level.Error(t.Logger).Log("msg", "Error read resp.body", "err", err)
			return err
		}
	} else {

		if t.KrbAuthMethod == "password" {
			data, err = lib.MakeKrb5RequestWithPassword(t.KrbPrincipal, t.KrbPassword, t.Url)

		} else if t.KrbAuthMethod == "keytab" {
			data, err = lib.MakeKrb5RequestWithKeytab(t.KrbKtPath, t.KrbPrincipal, t.Url)

		} else {
			level.Error(t.Logger).Log("msg", "Unsupported auth method")
			return fmt.Errorf("unsupported auth method")
		}

		if err != nil {
			level.Error(t.Logger).Log("msg", "Error make krb5 request", "err", err)
			return err
		}

	}

	t.BodyData = data

	var f interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		level.Error(t.Logger).Log("msg", "Error json Unmarshal", "err", err)
		return err
	}

	m := f.(map[string]interface{})
	// [{"name":"Hadoop:service=NameNode,name=FSNamesystem", ...}, {"name":"java.lang:type=MemoryPool,name=Code Cache", ...}, ...]
	var nameList = m["beans"].([]interface{})
	for _, nameData := range nameList {
		nameDataMap := nameData.(map[string]interface{})

		if strings.HasPrefix(nameDataMap["name"].(string), "Hadoop:service=") {

			regex := regexp.MustCompile(`Hadoop:service=(.*?),`)

			match := regex.FindStringSubmatch(nameDataMap["name"].(string))

			if len(match) > 1 {

				if match[1] == "HBase" {
					if strings.Contains(nameDataMap["name"].(string), "RegionServer") {
						t.ExporterName = "HbaseRegionServer"
						break
					} else if strings.Contains(nameDataMap["name"].(string), "Master") {
						t.ExporterName = "HbaseMaster"
						break
					}
				}

				t.ExporterName = match[1]
				break
			} else {
				level.Error(t.Logger).Log("msg", "Error pattern not match,unknown jmx service", "err", err)

				return fmt.Errorf("pattern not match,unknown jmx service")
			}
		}
	}
	return nil
}
