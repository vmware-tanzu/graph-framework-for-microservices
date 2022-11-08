from jira import ( JIRA,JIRAError )
import json
import requests
import os, sys

issue_types = [{"Bug":"bug"},{"Task":"feature request"}]
def parse_event_context():
    """
    Parse event context and run issue creation checks
    """
    try:
        create_issue_for = False
        event = os.environ.get("EVENT_CONTEXT")
        event_json = json.loads(event)
        issue_type = ""
        issue_body = event_json["issue"]["body"]
        issue_title = event_json["issue"]["title"]
        issue_url =  event_json["issue"]["url"]
        for label in event_json["issue"]["labels"]:
            for issue_type_obj in issue_types:
                for k,v in issue_type_obj.items():
                    if label["name"] ==  v:
                        issue_type = k
                        create_issue_for = True
        return issue_body, issue_type, issue_title, issue_url , create_issue_for
    except Exception as e:
        print("Please set EVENT_CONTEXT Variable with github issue event type... {}".format(e))
        sys.exit(1)


def run_create_issue(issue_body, issue_type, issue_title, issue_url):
    """
    JIRA_SERVER environment variable will be feed from Github CI for creating issues
    JIRA_USER, JIRA_PWD will be kept as a secret for authentication
    """
    JIRA_SERVER = os.environ.get("JIRA_SERVER")
    JIRA_USER = os.environ.get("JIRA_USER")
    JIRA_PWD = os.environ.get("JIRA_PWD")

    options = {"server": JIRA_SERVER}
    try:
        jira = JIRA(basic_auth=(JIRA_USER, JIRA_PWD), options=options)
    except JIRAError as e:
        print("Could not connect to JIRA due to Auth Failure/server not reachable: {}".format(e))
        sys.exit(1)
    except Exception as e:
        print("Could not connect to JIRA due to : {}".format(e))
        sys.exit(1)
    Body = "Github Issue: "+issue_url + "\n" + issue_body
    Issue = {
        'project': {'key': 'NPT'},
        'summary': issue_title ,
        'description': Body,
        'issuetype': {'name': issue_type },
    }
    try:
        new_issue = jira.create_issue(fields=Issue)
        Message="https://jira.eng.vmware.com/browse/"+new_issue.key
        URL=issue_url+"/comments"
        github_token = os.environ.get("GITHUBTOKEN")
        resp = requests.post(URL, json = {"body":"New JIRA Created with ID: {}".format(Message)},  headers={"Authorization":"Bearer {}".format(github_token),"Accept":"application/vnd.github+json"})
        print(resp.status_code)

    except Exception as e:
        print("Could not create JIRA due to : {}", e)
        sys.exit(1)

if __name__ == "__main__":
    """
    Run issue creation for bug and feature requests
    """
    issue_body, issue_type, issue_title, issue_url, create_issue_for = parse_event_context()
    if create_issue_for:
        run_create_issue(issue_body, issue_type, issue_title, issue_url)