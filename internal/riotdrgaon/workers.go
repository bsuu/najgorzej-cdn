package riotdrgaon

import (
	"fmt"
	"time"
)

func ScanVersion(v1 []string, v2 []string) (bool, []string) {
	versions := make([]string, 0)

	if len(v1) == len(v2) {
		return true, versions
	}

	for _, version := range v1 {
		same := false
		for _, localVersion := range v2 {
			if version == localVersion {
				same = true
			}
		}
		if !same {
			versions = append(versions, version)
		}
	}

	return len(versions) == 0, versions
}

func (r *RiotDragon) VersionsUpToDate() {
	riotVersions, err := r.GetGameVersions()

	if err != nil {
		return
	}

	upToDate, notFound := ScanVersion(riotVersions, r.VersionsIds)

	if !upToDate {
		r.VersionsIds = riotVersions
		for _, version := range notFound {
			_, err := r.DownloadVersion(version)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func VersionWorker(r *RiotDragon) {
	for {
		r.VersionsUpToDate()
		time.Sleep(15 * time.Minute)
	}
}
