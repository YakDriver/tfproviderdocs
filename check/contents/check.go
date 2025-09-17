package contents

import "fmt"

type CheckOptions struct {
	ArgumentsSection  *CheckArgumentsSectionOptions
	AttributesSection *CheckAttributesSectionOptions
	ExamplesSection   *CheckExamplesSectionOptions
	TimeoutsSection   *CheckTimeoutsSectionOptions
	ImportSection     *CheckImportSectionOptions

	DisallowAttributesSection          bool
	AttributesSectionDisallowedMessage string
	DisallowImportSection              bool
	ImportSectionDisallowedMessage     string
}

func (d *Document) Check(opts *CheckOptions) error {
	d.CheckOptions = opts

	if err := d.checkTitleSection(); err != nil {
		return err
	}

	if err := d.checkExampleSection(); err != nil {
		return err
	}

	if err := d.checkArgumentsSection(); err != nil {
		return err
	}

	if d.CheckOptions != nil && d.CheckOptions.DisallowAttributesSection {
		if d.Sections.Attributes != nil {
			msg := "attribute section is not allowed"

			if d.CheckOptions.AttributesSectionDisallowedMessage != "" {
				msg = d.CheckOptions.AttributesSectionDisallowedMessage
			}

			return fmt.Errorf(msg)
		}
	} else {
		if err := d.checkAttributesSection(); err != nil {
			return err
		}
	}

	if err := d.checkTimeoutsSection(); err != nil {
		return err
	}

	if d.CheckOptions != nil && d.CheckOptions.DisallowImportSection {
		if d.Sections.Import != nil {
			msg := "import section is not allowed"

			if d.CheckOptions.ImportSectionDisallowedMessage != "" {
				msg = d.CheckOptions.ImportSectionDisallowedMessage
			}

			return fmt.Errorf(msg)
		}
	} else {
		if err := d.checkImportSection(); err != nil {
			return err
		}
	}

	return nil
}
