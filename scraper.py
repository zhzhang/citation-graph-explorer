import urllib.request
import urllib.error
import json
import pickle
import random

seen = set()
stack = [('arXiv:1701.01821', 0)]
output = []
sampling_dampener = 0.01

while len(stack) > 0:
    print(len(stack))
    paper_id, depth = stack.pop()
    if paper_id in seen:
        continue
    url = "http://api.semanticscholar.org/v1/paper/%s" % paper_id
    try:
        contents = urllib.request.urlopen(url).read()
    except urllib.error.HTTPError:
        continue
    seen.add(paper_id)
    contents = json.loads(contents)
    output.append(contents)
    for paper in contents['citations']:
        if random.random() < ((sampling_dampener / (depth + 1)) ** depth):
            stack.append((paper['paperId'], depth + 1))
    for paper in contents['references']:
        if random.random() < ((sampling_dampener / (depth + 1)) ** depth):
            stack.append((paper['paperId'], depth + 1))


with open('samples.pkl', 'wb') as f:
    pickle.dump(output, f)
