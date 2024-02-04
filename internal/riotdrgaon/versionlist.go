package riotdrgaon

type VersionList struct {
	size     int
	versions []*Version
}

func (v *VersionList) Get(versionId string) *Version {
	for _, version := range v.versions {
		if version.Id == versionId {
			return version
		}
	}
	return nil
}

func (v *VersionList) Containes(versionId string) (bool, int) {
	for i, version := range v.versions {
		if version.Id == versionId {
			return true, i
		}
	}
	return false, -1
}

func (v *VersionList) Put(version *Version) {

	//If the version is already in the list, we remove it
	if contains, index := v.Containes(version.Id); contains {
		v.versions = append(v.versions[:index], v.versions[index+1:]...)
	}

	//If the list is empty, we add the version and return
	if len(v.versions) == 0 {
		v.versions = make([]*Version, v.size)
		v.versions[0] = version
		return
	}

	//If the list is full, we remove the last element
	if len(v.versions) == v.size {
		v.versions = v.versions[:len(v.versions)-1]
	}

	//Shifting by one down
	for i := len(v.versions) - 1; i > 0; i-- {
		v.versions[i] = v.versions[i-1]
	}
	//Setting the first element to the new version
	v.versions[0] = version
}
