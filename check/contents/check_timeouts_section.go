package contents

import "fmt"

type CheckTimeoutsSectionOptions struct {
	RequireSection SectionRequirement
}

func (d *Document) checkTimeoutsSection() error {
	checkOpts := &CheckTimeoutsSectionOptions{}

	if d.CheckOptions != nil && d.CheckOptions.TimeoutsSection != nil {
		checkOpts = d.CheckOptions.TimeoutsSection
	}

	section := d.Sections.Timeouts

	if section == nil {
		return nil
	} else {
		if checkOpts.RequireSection == Forbidden {
			return fmt.Errorf("timeouts section should not be present")
		}
	}

	return nil
}
