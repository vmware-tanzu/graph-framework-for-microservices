Thank you for investing your time in contributing to Nexus SDK ✨.

# Contributor Guide

## Merge Requests / Pull Requests

### Commit Messages

Commit messages on the Merge Request should be of the format:

> **NPT-ABCD <[Summary](CONTRIBUTING.md#summary)>**
> 
> **< blank line >**
> 
> **[Problem statement](CONTRIBUTING.md#problem-statement) - can be description of bug, reason for change,
> need, ask etc.**
> 
> **< blank line >**
> 
> **[Description of fix](CONTRIBUTING.md#description-of-fix) / code change.**




#### Summary

* A properly formed git commit subject line should always be able to complete the following sentence
  If applied, this commit will <your subject line here>
* Should state with Jira ID. Jira ID whould always be capitalized.
* Do not end the subject line with a period
* Use the imperative mood in the subject line
#### Problem statement
* Separate problem statement from subject, with a blank line
* Describe why a change is being made.
* Bullet points are okay.
* Do not assume the reviewer understands what the original problem was.
#### Description of fix
* Separate explanation from problem statement, with a blank line
* Bullet points are okay, too.
* Leave out details about how a change has been made, unless it is necessary for clarity.

An example / simple commit message, based on above quidelines.

```
NPT-239 Add a FAQ page to document important / frequent questions
    
In our documentation, there needs to be a place for info that come
up a questions frequently. The workflow page is sometime too verbose
to answer such questions quickly.
    
Add a FAQ page to the docs, where such questions can be captured.
```
