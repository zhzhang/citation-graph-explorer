import urllib.request
import urllib.error
import json
import random

f = open('output.json', 'wb')
seen = set()
stack = [('arXiv:1701.01821', 0)]

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
    f.write(contents)
    contents = json.loads(contents)
    for paper in contents['citations']:
        if random.random() < ((0.2 / (depth + 1)) ** depth):
            stack.append((paper['paperId'], depth + 1))
    for paper in contents['references']:
        if random.random() < ((0.2 / (depth + 1)) ** depth):
            stack.append((paper['paperId'], depth + 1))

f.close()
