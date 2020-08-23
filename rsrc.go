// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package rsrc

import (
	"encoding/binary"
	"io"
	"os"
)

type Resource struct {
	Type       string
	ID         int16
	Name       string
	Data       []byte
	dataOffset int64
}

type ResourceFile struct {
	file            io.ReadSeeker
	resourcesByType map[string][]Resource
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
	if _, err := rf.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

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
	numberOfTypes++

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
		numberOfResources++

		var referenceListOffset uint16
		if err := binary.Read(rf.file, binary.BigEndian, &referenceListOffset); err != nil {
			return err
		}

		rf.resourcesByType[fourCharacterCode(resourceType)] = make([]Resource, 0)

		// Read the resource reference list

		for referenceIndex := 0; referenceIndex < int(numberOfResources); referenceIndex++ {
			// TODO Note how that +2 offset is now missing?
			offset := int64(resourceMapOffset) + int64(resourceTypeListOffset) + int64(referenceListOffset) + int64(referenceIndex)*12
			if _, err := rf.file.Seek(offset, 0); err != nil {
				return err
			}

			var resourceID int16
			if err := binary.Read(rf.file, binary.BigEndian, &resourceID); err != nil {
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

			// Read the name

			var resourceName string

			if resourceNameOffset != -1 {
				if _, err := rf.file.Seek(int64(resourceMapOffset)+int64(resourceNameListOffset)+int64(resourceNameOffset), 0); err != nil {
					return err
				}

				name, err := readPascalString(rf.file)
				if err != nil {
					return err
				}
				resourceName = name
			}

			// We have now fully read one resource

			resource := Resource{
				Type:       fourCharacterCode(resourceType),
				ID:         resourceID,
				Name:       resourceName,
				dataOffset: int64(resourceDataOffset) + int64(thisResourceDataOffset),
			}

			rf.resourcesByType[fourCharacterCode(resourceType)] = append(rf.resourcesByType[fourCharacterCode(resourceType)], resource)
		}
	}

	return nil
}

func New(rs io.ReadSeeker) (*ResourceFile, error) {
	rf := &ResourceFile{
		file:            rs,
		resourcesByType: make(map[string][]Resource),
	}
	if err := rf.parseResourceMap(); err != nil {
		return nil, err
	}
	return rf, nil
}

func FromPath(path string) (*ResourceFile, error) {
	file, err := os.Open(path + "/..namedfork/rsrc")
	if err != nil {
		return nil, err
	}
	return New(file)
}
