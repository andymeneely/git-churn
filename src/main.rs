use git2::*;
use std::path::Path;

mod git_churn;

fn compute_churn() {
    let repo = Repository::open(".").expect("Could not open repo");
    let head = repo
        .revparse_single("test-first-test")
        .expect("Failed revparse");
    let commit = head.peel_to_commit().expect("Could not peel to commit");
    let commit_tree = commit.tree().expect("Could not get tree");
    for p in commit.parents() {
        let parent_tree = p.tree().expect("Failed find parent tree");
        let mut diff_opts = init_diff_opts();
        let diff = repo
            .diff_tree_to_tree(Some(&parent_tree), Some(&commit_tree), Some(&mut diff_opts))
            .unwrap();
        let diff_stats = diff.stats().expect("Failed to compute diff_stats");

        println!(
            "Diff stats: {} insertions, {} deletions",
            diff_stats.insertions(),
            diff_stats.deletions()
        );

        println!("--- DIFF ---");
        diff.foreach(
            &mut file_cb,
            Some(&mut binary_cb),
            Some(&mut hunk_cb),
            Some(&mut line_cb),
        )
        .expect("Failed to iterate over diff");

        println!("-- BLAME --");
        let mut blame_opts = init_blame_opts(head.id());
        let path = Path::new("src/main.rs");
        let blame = repo
            .blame_file(&path, Some(&mut blame_opts))
            .expect("Failed to execute blame");
        for hunk in blame.iter() {
            println!(
                "Hunk {} Start line: {}",
                hunk.orig_commit_id(),
                hunk.final_start_line()
            );
        }
    }
}

fn init_diff_opts() -> DiffOptions {
    let mut diff_opts = DiffOptions::new();
    diff_opts.ignore_whitespace(true);
    diff_opts.context_lines(0);
    diff_opts.ignore_filemode(true);
    diff_opts.indent_heuristic(true);
    return diff_opts;
}

fn init_blame_opts(head_id: Oid) -> BlameOptions {
    let mut blame_opts = BlameOptions::new();
    blame_opts.newest_commit(head_id);
    blame_opts.track_copies_same_file(true);
    blame_opts.track_copies_same_commit_moves(false);
    blame_opts.track_copies_same_commit_copies(false);
    blame_opts.track_copies_any_commit_copies(false);
    return blame_opts;
}

fn main() {
    compute_churn();
}

fn file_cb(_d: DiffDelta, _progress: f32) -> bool {
    true
}

fn binary_cb(_d: DiffDelta, _db: DiffBinary) -> bool {
    true
}

fn hunk_cb(_d: DiffDelta, _dh: DiffHunk) -> bool {
    true
}

fn line_cb(_d: DiffDelta, _o: Option<DiffHunk>, l: DiffLine) -> bool {
    print!(
        "{}   {}",
        l.origin(),
        String::from_utf8(l.content().to_vec()).unwrap()
    );
    true
}

#[cfg(test)]
mod tests {

    #[test]
    fn test_goofing_tag() {
        let repo = git2::Repository::open(".").expect("Failed to open git repo .");
        let head = repo
            .revparse_single("test-goofing")
            .expect("Could not parse");
        assert_eq!(
            "79caa008ba1f9d06b34b4acc7c03d7fade185a63",
            head.id().to_string()
        );
    }
}
