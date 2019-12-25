import argparse
import pickle


KEYS = [
    "authors",
    "citationVelocity",
    "citations",
    "influentialCitationCount",
    "paperId",
    "title",
    "url",
    "venue",
    "year",
]


def parse_record(record):
    output = {k: record[k] for k in KEYS}
    return output


def load_articles(path):
    f = open(path, 'rb')
    records = pickle.loads(f.read())
    output = {r["paperId"]: r for r in records}
    for r in records:
        if "citations" in r.keys():
            citations = r["citations"]
            for c in citations:
                if not c["paperId"] in output:
                    c["citations"] = []
                    output[c["paperId"]] = c
    return output


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Load the citation graph into memory.")
    parser.add_argument("path")
    args = parser.parse_args()
    path = args.path
    load_articles(path)
