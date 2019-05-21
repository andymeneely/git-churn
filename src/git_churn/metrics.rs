use git2::*;
use std::path::Path;
use super::*;

pub fn compute_churn(repo: &Repository, commit: &Commit, _interactive: bool) -> Stats {
    let commit_tree = commit.tree().expect("Could not get tree");
    let stats = Stats::new();
    for p in commit.parents() {
        let parent_tree = p.tree().expect("Failed find parent tree");
        let mut diff_opts = init_diff_opts();
        let diff = repo
            .diff_tree_to_tree(Some(&parent_tree), Some(&commit_tree), Some(&mut diff_opts))
            .expect("Failed to diff tree to tree");
        // let diff_stats = diff.stats().expect("Failed to compute diff_stats");
        // stats.set(diff_stats);

        println!("--- DIFF ---");
        diff.foreach(&mut file_cb, None, None, Some(&mut line_cb))
            .expect("Failed to iterate over diff");

        println!("-- BLAME --");
        let mut blame_opts = init_blame_opts(commit.id());
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
    return stats;
}

fn init_diff_opts() -> DiffOptions {
    let mut diff_opts = DiffOptions::new();
    diff_opts.ignore_whitespace(true);
    diff_opts.context_lines(0);
    diff_opts.ignore_filemode(true);
    diff_opts.indent_heuristic(true);
    return diff_opts;
}


fn file_cb(_d: DiffDelta, _progress: f32) -> bool {
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

fn init_blame_opts(head_id: Oid) -> BlameOptions {
    let mut blame_opts = BlameOptions::new();
    blame_opts.newest_commit(head_id);
    blame_opts.track_copies_same_file(true);
    blame_opts.track_copies_same_commit_moves(false);
    blame_opts.track_copies_same_commit_copies(false);
    blame_opts.track_copies_any_commit_copies(false);
    return blame_opts;
}
