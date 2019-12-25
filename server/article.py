from .utils.load_articles import load_articles
import json


SAMPLE_PATH = "./samples.pkl"


class ArticleResource:
    def __init__(self):
        self.articles = load_articles(SAMPLE_PATH)
        print("READY")

    def on_get(self, req, resp):
        paper_id = req.get_param("paperId")
        depth = int(req.get_param("depth"))
        if depth is None:
            resp.body = json.dumps(self.articles[paper_id])
        else:
            resp.body = json.dumps(self._bfs(paper_id, depth))

    def _bfs(self, root_id, max_depth):
        fringe = [(root_id, 0)]
        results_by_depth = {d: (set(), []) for d in range(max_depth + 1)}
        while len(fringe) > 0:
            paper_id, depth = fringe.pop(0)
            paper = self.articles[paper_id]
            seen, output = results_by_depth[depth]
            if paper_id not in seen:
                seen.add(paper_id)
                output.append(paper)
            if not depth == max_depth:
                for c in paper["citations"]:
                    fringe.append((c["paperId"], depth + 1))
        return {d: v[1] for d, v in results_by_depth.items()}
