---
name: Bug report
about: Create a report to help us improve
title: ''
labels: bug
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior.

**Expected behavior**
A clear and concise description of what you expected to happen.

**Version:**
 - Output of the command:    nexus version

**Debug:**

1.  Output of the command:
      - nexus debug
     
3.  For build failures:                  
      - nexus datamodel build --debug

4.  For installation failures:
      - nexus prereq verify
      -  kubectl get pods -A -o yaml


**System:**
 - OS: [e.g. iOS]

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Additional context**
Add any other context about the problem here.

**Requester**
Select from the list of well known requesters or select "Community"
- [x] Community
- [ ] Project ServiceMesh
- [ ] Project Mazinger
- [ ] Project Watch
