use super::*;
// use git2::*;
// use std::path::Path;

pub fn compute_churn(repo: &Repository, commit: &Commit, _interactive: bool) -> Stats {
    let commit_tree = commit.tree().expect("Could not get tree");
    let commit_str: String = commit.id().to_string();

    let mut commit_stats = CommitStats::new();
    commit_stats.merge_commit = commit.parents().count() > 1;
    for p in commit.parents() {

        println!("--- DIFF ---");
        let parent_tree = p.tree().expect("Failed find parent tree");
        let mut diff_opts = init_diff_opts();
        let diff = repo
            .diff_tree_to_tree(Some(&parent_tree), Some(&commit_tree), Some(&mut diff_opts))
            .expect("Failed to diff tree to tree");
        // let diff_stats = diff.stats().expect("Failed to compute diff_stats");
        // stats.set(diff_stats);
        // let file_edits:HashMap<String, CommitPathStats> = HashMap::new();

        let mut file_cb = |d: DiffDelta, _progress: f32| -> bool {
            let path_str = path_from(d);
            println!("Processing file: {}", path_str);
            // commit_stats.commit_path_stats.insert(path_str, CommitPathStats::new());
            return true;
        };

        let mut line_cb = |d: DiffDelta, _o: Option<DiffHunk>, l: DiffLine| -> bool {
            // let path_str = path_from(d);
            //
            // if !commit_stats.commit_path_stats.contains_key(&path_str) {
            //     commit_stats.commit_path_stats.insert(path_from(d), CommitPathStats::new());
            // }
            // let mut cps = commit_stats.commit_path_stats.get_mut(&path_str).expect("Commit path stats not initialized");

            // match l.origin() {
            //     '+' => cps.insertions += 1,
            //     '-' => cps.deletions += 1,
            //     _ => (),
            // }
            print!(
                "{}   {}",
                l.origin(),
                String::from_utf8(l.content().to_vec()).unwrap()
            );
            true
        };

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
    let mut stats = Stats::new();
    stats.commit_stats.insert(commit_str, commit_stats);
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

fn init_blame_opts(head_id: Oid) -> BlameOptions {
    let mut blame_opts = BlameOptions::new();
    blame_opts.newest_commit(head_id);
    blame_opts.track_copies_same_file(true);
    blame_opts.track_copies_same_commit_moves(false);
    blame_opts.track_copies_same_commit_copies(false);
    blame_opts.track_copies_any_commit_copies(false);
    return blame_opts;
}

fn path_from(d: DiffDelta) -> String {
    let path = d
        .new_file()
        .path()
        .expect("Could not find new file in delta");
    return String::from(path.to_str().expect("Could not convert to path string."));
}
