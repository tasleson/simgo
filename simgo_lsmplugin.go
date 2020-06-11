package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	lsm "github.com/libstorage/libstoragemgmt-golang"
	"github.com/libstorage/libstoragemgmt-golang/errors"
)

type pluginData struct {
	reg lsm.PluginRegister
	c   *lsm.ClientConnection
	tmo uint32
}

var state pluginData

func register(p *lsm.PluginRegister) error {

	parsed, pE := url.Parse(p.URI)
	if pE != nil {
		return pE
	}

	values := strings.Split(parsed.RawQuery, "=")
	if len(values) != 2 || values[0] != "forward" {
		return &errors.LsmError{
			Code:    errors.InvalidArgument,
			Message: fmt.Sprintf("expected query string to be 'forward=<otherplugin>' got %s", parsed.RawQuery)}
	}

	var cE error
	state.c, cE = lsm.Client(fmt.Sprintf("%s://", values[1]), "", p.Timeout)
	return cE
}

func unregister() error {
	return state.c.Close()
}

func systems() ([]lsm.System, error) {
	return state.c.Systems()
}

func pools(search ...string) ([]lsm.Pool, error) {
	if len(search) > 0 {
		return state.c.Pools(search[0], search[1])
	}
	return state.c.Pools()
}

func volumes(search ...string) ([]lsm.Volume, error) {
	if len(search) > 0 {
		return state.c.Volumes(search[0], search[1])
	}
	return state.c.Volumes()
}

func tmoSet(timeout uint32) error {
	state.tmo = timeout
	return nil
}

func tmoGet() uint32 {
	return state.tmo
}

func capabilities(system *lsm.System) (*lsm.Capabilities, error) {
	return state.c.Capabilities(system)
}

func jobStatus(jobID string) (*lsm.JobInfo, error) {
	var item interface{}
	jobStatus, jobPercent, error := state.c.JobStatus(jobID, &item)
	if error != nil {
		return nil, error
	}
	return &lsm.JobInfo{Status: jobStatus, Percent: jobPercent, Item: item}, nil
}

func jobFree(jobID string) error {
	return state.c.JobFree(jobID)
}

func volCreate(pool *lsm.Pool, volumeName string, size uint64,
	provisioning lsm.VolumeProvisionType) (*lsm.Volume, *string, error) {

	var volume lsm.Volume
	jobID, error := state.c.VolumeCreate(pool, volumeName, size, provisioning,
		false, &volume)

	if jobID != nil {
		return nil, jobID, error
	}
	return &volume, nil, error
}

func volDelete(volume *lsm.Volume) (*string, error) {
	return state.c.VolumeDelete(volume, false)
}

func disks() ([]lsm.Disk, error) {
	return state.c.Disks()
}

func volReplicate(optionalPool *lsm.Pool, repType lsm.VolumeReplicateType,
	sourceVolume *lsm.Volume, name string) (*lsm.Volume, *string, error) {

	var volume lsm.Volume
	jobID, error := state.c.VolumeReplicate(optionalPool, repType, sourceVolume, name, false, &volume)

	if jobID != nil {
		return nil, jobID, error
	}
	return &volume, nil, error
}

func volReplicateRange(repType lsm.VolumeReplicateType, srcVol *lsm.Volume, dstVol *lsm.Volume,
	ranges []lsm.BlockRange) (*string, error) {
	return state.c.VolumeReplicateRange(repType, srcVol, dstVol, ranges, false)
}

func volRepRangeBlockSize(system *lsm.System) (uint32, error) {
	return state.c.VolumeRepRangeBlkSize(system)
}

func volResize(vol *lsm.Volume, newSizeBytes uint64) (*lsm.Volume, *string, error) {

	var volume lsm.Volume
	jobID, error := state.c.VolumeResize(vol, newSizeBytes, false, &volume)
	if jobID != nil {
		return nil, jobID, error
	}
	return &volume, nil, error
}

func main() {
	var cb lsm.CallBacks
	cb.Required.Systems = systems
	cb.Required.PluginRegister = register
	cb.Required.PluginUnregister = unregister
	cb.Required.Pools = pools
	cb.Required.TimeOutSet = tmoSet
	cb.Required.TimeOutGet = tmoGet
	cb.Required.Capabilities = capabilities
	cb.Required.JobStatus = jobStatus
	cb.Required.JobFree = jobFree

	cb.San.VolumeCreate = volCreate
	cb.San.VolumeDelete = volDelete
	cb.San.Volumes = volumes
	cb.San.Disks = disks
	cb.San.VolumeReplicate = volReplicate
	cb.San.VolumeReplicateRange = volReplicateRange
	cb.San.VolumeRepRangeBlkSize = volRepRangeBlockSize
	cb.San.VolumeResize = volResize

	plugin, err := lsm.PluginInit(&cb, os.Args, "golang forwarding plugin", "0.0.1")
	if err != nil {
		fmt.Printf("Failed to initialize plugin, exiting! (%s)\n", err)
	} else {
		plugin.Run()
	}
}
