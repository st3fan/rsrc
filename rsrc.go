// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package rsrc

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

func fourCharacterCode(code uint32) string {
	s := ""
	s += string(code >> 24 & 0x000000ff)
	s += string(code >> 16 & 0x000000ff)
	s += string(code >> 8 & 0x000000ff)
	s += string(code & 0x000000ff)
	return s
}

func readPascalString(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)

	length, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	name := ""
	for i := 0; i < int(length); i++ {
		c, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		name += string(c)
	}

	return name, nil
}

type Resource struct {
	Type       string
	ID         int16
	Name       string
	Data       []byte
	dataOffset int64
}

type ResourceFile struct {
	file            *os.File
	resourcesByType map[string][]Resource
}

func (rf *ResourceFile) Close() error {
	return rf.file.Close()
}

func (rf *ResourceFile) CountResources(resourceType string) int {
	return len(rf.resourcesByType[resourceType])
}

func (rf *ResourceFile) GetResource(resourceType string, resourceIndex int) (Resource, bool) {
	if resourceIndex < len(rf.resourcesByType[resourceType]) {
		resource := rf.resourcesByType[resourceType][resourceIndex]

		if _, err := rf.file.Seek(resource.dataOffset, 0); err != nil {
			return Resource{}, false
		}

		var resourceLength int32
		if err := binary.Read(rf.file, binary.BigEndian, &resourceLength); err != nil {
			return Resource{}, false
		}

		resource.Data = make([]byte, resourceLength)
		if _, err := rf.file.Read(resource.Data); err != nil {
			return Resource{}, false
		}

		return resource, true
	}
	return Resource{}, false
}

func (rf *ResourceFile) parseResourceMap() error {
	// Read the header

	var resourceDataOffset uint32
	if err := binary.Read(rf.file, binary.BigEndian, &resourceDataOffset); err != nil {
		return err
	}

	var resourceMapOffset uint32
	if err := binary.Read(rf.file, binary.BigEndian, &resourceMapOffset); err != nil {
		return err
	}

	var resourceDataLength uint32
	if err := binary.Read(rf.file, binary.BigEndian, &resourceDataLength); err != nil {
		return err
	}

	var resourceMapLength uint32
	if err := binary.Read(rf.file, binary.BigEndian, &resourceMapLength); err != nil {
		return err
	}

	// Read the resource map

	if _, err := rf.file.Seek(int64(resourceMapOffset)+24, 0); err != nil {
		return err
	}

	var resourceTypeListOffset uint16
	if err := binary.Read(rf.file, binary.BigEndian, &resourceTypeListOffset); err != nil {
		return err
	}

	var resourceNameListOffset uint16
	if err := binary.Read(rf.file, binary.BigEndian, &resourceNameListOffset); err != nil {
		return err
	}

	var numberOfTypes uint16
	if err := binary.Read(rf.file, binary.BigEndian, &numberOfTypes); err != nil {
		return err
	}
	numberOfTypes += 1

	//fmt.Printf("resourceTypeListOffset = %v\n", resourceTypeListOffset)
	//fmt.Printf("resourceNameListOffset = %v\n", resourceNameListOffset)
	//fmt.Printf("numberOfTypes = %v\n", numberOfTypes)

	// Read the resource type list

	// TODO Where does that +2 offset come from. Does not match Inside Macintosh.

	for typeIndex := 0; typeIndex < int(numberOfTypes); typeIndex++ {
		if _, err := rf.file.Seek(int64(resourceMapOffset)+int64(resourceTypeListOffset)+2+int64(typeIndex)*8, 0); err != nil {
			return err
		}

		var resourceType uint32
		if err := binary.Read(rf.file, binary.BigEndian, &resourceType); err != nil {
			return err
		}

		var numberOfResources uint16
		if err := binary.Read(rf.file, binary.BigEndian, &numberOfResources); err != nil {
			return err
		}
		numberOfResources += 1

		var referenceListOffset uint16
		if err := binary.Read(rf.file, binary.BigEndian, &referenceListOffset); err != nil {
			return err
		}

		//fmt.Printf("   resourceType = %s\n", fourCharacterCode(resourceType))
		//fmt.Printf("   numberOfResources = %v\n", numberOfResources)
		//fmt.Printf("   referenceListOffset = %v\n", referenceListOffset)

		rf.resourcesByType[fourCharacterCode(resourceType)] = make([]Resource, 0)

		// Read the resource reference list

		for referenceIndex := 0; referenceIndex < int(numberOfResources); referenceIndex++ {
			// TODO Note how that +2 offset is now missing?
			offset := int64(resourceMapOffset) + int64(resourceTypeListOffset) + int64(referenceListOffset) + int64(referenceIndex)*12
			if _, err := rf.file.Seek(offset, 0); err != nil {
				return err
			}

			var resourceId int16
			if err := binary.Read(rf.file, binary.BigEndian, &resourceId); err != nil {
				return err
			}

			var resourceNameOffset int16
			if err := binary.Read(rf.file, binary.BigEndian, &resourceNameOffset); err != nil {
				return err
			}

			var thisResourceDataOffset int32
			if err := binary.Read(rf.file, binary.BigEndian, &thisResourceDataOffset); err != nil {
				return err
			}
			thisResourceDataOffset &= 0x00ffffff

			//fmt.Printf("      resourceId = %v\n", resourceId)
			//fmt.Printf("      resourceNameOffset = %v\n", resourceNameOffset)
			//fmt.Printf("      resourceDataOffset = %v\n", resourceDataOffset)

			// Read the name

			var resourceName string = ""

			if resourceNameOffset != -1 {
				if _, err := rf.file.Seek(int64(resourceMapOffset)+int64(resourceNameListOffset)+int64(resourceNameOffset), 0); err != nil {
					return err
				}

				name, err := readPascalString(rf.file)
				if err != nil {
					return err
				}
				resourceName = name

				//fmt.Printf("         resourceName = %s\n", resourceName)
			}

			// We have now fully read one resource

			resource := Resource{
				Type:       fourCharacterCode(resourceType),
				ID:         resourceId,
				Name:       resourceName,
				dataOffset: int64(resourceDataOffset) + int64(thisResourceDataOffset),
			}

			rf.resourcesByType[fourCharacterCode(resourceType)] = append(rf.resourcesByType[fourCharacterCode(resourceType)], resource)
		}
	}

	return nil
}

func Open(path string) (*ResourceFile, error) {
	file, err := os.Open(path + "/..namedfork/rsrc")
	if err != nil {
		return nil, err
	}
	rf := &ResourceFile{
		file:            file,
		resourcesByType: make(map[string][]Resource),
	}

	if err := rf.parseResourceMap(); err != nil {
		rf.Close()
		return nil, err
	}

	return rf, nil
}