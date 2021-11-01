// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package network

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Links struct {
	LinksByMAC map[string]Link
}

type Link struct {
	Name      string
	Ifindex   int
	OperState string
	Mac       string
	MTU       int
	Addresses *map[string]bool
}

func AcquireLinks() (Links, error) {
	linkList, err := netlink.LinkList()
	if err != nil {
		return Links{}, err
	}

	links := make(map[string]Link)
	for _, link := range linkList {
		if link.Attrs().Name == "lo" {
			continue
		}

		l := Link{
			Name:      link.Attrs().Name,
			Ifindex:   link.Attrs().Index,
			Mac:       link.Attrs().HardwareAddr.String(),
			OperState: link.Attrs().OperState.String(),
			MTU:       link.Attrs().MTU,
		}

		links[link.Attrs().HardwareAddr.String()] = l

		log.Debugf("Acquired link='%s' ifindex='%d'", link.Attrs().Name, link.Attrs().Index)
	}

	return Links{
		LinksByMAC: links,
	}, nil
}

func GetLinkMacByIndex(links *Links, ifIndex int) (string, error) {
	for _, l := range links.LinksByMAC {
		if l.Ifindex == ifIndex {
			return l.Mac, nil
		}
	}

	return "", errors.New("not found")
}

func LinkSetOperStateUp(ifIndex int) error {
	link, err := netlink.LinkByIndex(ifIndex)
	if err != nil {
		return err
	}

	err = netlink.LinkSetUp(link)
	if err != nil {
		return err
	}

	return nil
}

func LinkSetMtu(ifIndex int, mtu int) error {
	link, err := netlink.LinkByIndex(ifIndex)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetMTU(link, mtu); err != nil {
		return err
	}

	return nil
}
