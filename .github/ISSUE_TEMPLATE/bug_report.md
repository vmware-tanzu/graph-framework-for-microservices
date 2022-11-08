name: Bug Report
description: File a bug report
title: "[Bug]: "
labels: ["bug"]
assignees:
- xmen4xp
  body:
- type: markdown
  attributes:
  value: |
  Thanks for taking the time to fill out this bug report!
- type: input
  id: contact
  attributes:
    label: Contact Details
    description: How can we get in touch with you if we need more info?
    placeholder: ex. email@example.com
    validations:
    required: false
- type: textarea
    id: version
    attributes:
    label: What version are you running?
    description: Output of: nexus version
    placeholder: Output of: nexus version
    value: "Output of: nexus version"
    validations:
    required: true
- type: textarea
  id: what-happened
  attributes:
    label: What happened?
    description: Describe the bug
    placeholder: Tell us what you see!
    value: "A bug happened!"
    validations:
    required: true
- type: textarea
  id: what-was-expected-behavior
  attributes:
    label: Describe the expected behavior
    description: What was the expected behavior?
    placeholder: What was the expected behavior?
    value: "Not Applicable"
    validations:
    required: false
- type: dropdown
  id: priority
  attributes:
  label: How critical is this bug to you?
  multiple: false
  options:
  - Blocker - solution is unusable without this feature
  - Critical - solution is severely limited in value without this feature
  - Major - important feature to be incorporated as there are no  known alternatives in the solution
  - Minor - nice to have feature that adds value to the solution
- type: textarea
  id: how-can-we-recreate-the-bug
  attributes:
    label: How can we recreate the bug?
    description: Share with us, the steps to hit the bug
    placeholder: Share with us, the steps to hit the bug
    value: "Not Applicable"
    validations:
    required: false
- type: markdown
    attributes:
    value: |
    Providing additional details would speed up resolution of the issue by many light years !
    For generic debug info, attach output of:
    '''
    nexus debug
    '''
  
    For compilation bugs, output of:
    '''
    nexus datamodel build --debug
    '''

    For installation bugs, output of:
    '''
    nexus prereq verify
    kubectl get pods -A -o yaml
    '''
- type: textarea
  id: debug
  attributes:
  label: Any debug data that you are able to share?
  description: Any debug data that you are able to share?
  placeholder: Any debug data that you are able to share?
  value: "Not Applicable"
  validations:
  required: false
- type: dropdown
  id: associated-project
  attributes:
  label: Tell us the project / group you are associated with
  description: Tell us the project / group you are associated with
  options:
  - Community (Default)
  - Project ServiceMesh
  - Project Mazinger
  - Project Watch
  validations:
  required: true
- type: dropdown
  id: operating system
  attributes:
  label: What is your operating system?
  multiple: false
  options:
  - Linux
  - MacOS
- type: textarea
  id: additional-info
  attributes:
  label: Any additional / relevant info
  description: Any additional / relevant info.
  value: "Not Applicable"
  validations:
  required: false
