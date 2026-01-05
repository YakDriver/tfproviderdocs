// Copyright (c) IBM Corp. 2019-2026
// SPDX-License-Identifier: MPL-2.0

// Copyright (c) HashiCorp, Inc. 2019-2026
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
		if checkOpts.RequireSection == Required {
			return fmt.Errorf("missing timeouts section: ## Timeouts")
		}
		return nil
	} else {
		if checkOpts.RequireSection == Forbidden {
			return fmt.Errorf("timeouts section should not be present")
		}
	}

	return nil
}
