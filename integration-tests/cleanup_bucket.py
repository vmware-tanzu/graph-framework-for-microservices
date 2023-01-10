import sys, subprocess, re
from datetime import datetime , date

from numpy import diff


def get_folderStats(folder):
    statsObject = []
    print("Combing through {}".format(folder))
    output = subprocess.Popen(["gsutil","ls","-l","gs://{}/".format(folder) ],  stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in output.stdout.readlines():
        result = line.decode("utf-8").strip()
        statsObject.append(result)
    return statsObject

def deleteFolder(folder):
    return subprocess.call(["gsutil","rm","-r","gs://{}/".format(folder) ] )

VersionRegex = re.compile(r'v?\d+\.\d+\.\d+.*')
if __name__ == "__main__":
    bucket_name = sys.argv[1]
    if len(sys.argv) > 2:
        days = sys.argv[2]
    else:
        days = 10
    if bucket_name == "":
       print("Please provide bucket name to scoop through...")
       exit(1)
    output = get_folderStats(bucket_name)
    for value in output: 
        if not value.endswith("tar") and value.startswith("gs:"):
           folderName = value.split("gs://")[1].split("/")[-2]
           if not re.match(VersionRegex, folderName):
              folderOutput = get_folderStats("{}/{}".format(bucket_name, folderName))
              for files in folderOutput:
                  if "gs://" in files:
                      try:
                        dateToParse = files.split("  ")[1]
                        dateObj = datetime.strptime(dateToParse, "%Y-%m-%dT%H:%M:%SZ")
                        currDate = datetime.now()
                        difference = (currDate - dateObj)
                        if difference.days > int(days):
                            print("{} created before days: {}".format(folderName, days)) 
                            try:
                                if len(sys.argv) > 3:
                                    if sys.argv[3] == "delete":
                                        print("Deleting the folder which is a commit ID and more than 10 days :{}".format(folderName)) 
                                        deleteFolder(folderName) 
                            except Exception as e:
                                print("error in arguments check from CLI due to {}".format(e))
                      except Exception as e:
                        print("could not comb through {} due to {}".format(folderName, e))
                        pass
             