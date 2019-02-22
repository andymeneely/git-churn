use git2::*;

#[allow(unused_variables)]
fn main() {
    let repo = match Repository::open(".") {
        Ok(repo) => repo,
        Err(e) => panic!("failed to open: {}", e),
    };
    let head = match repo.revparse_single("test-first-test") {
        Ok(head) => head,
        Err(e) => panic!("failed to revparse: {}", e),
    };
    let commit  = head.peel_to_commit().unwrap();
    let commit_tree = match head.peel_to_commit().unwrap().tree() {
        Ok(commit_tree) => commit_tree,
        Err(e) => panic!("failed to revparse: {}", e),
    };
    println!("Commit message is: {}",commit.message().unwrap_or(""));
    for p in commit.parents(){
        let p_tree = match p.tree() {
            Ok(p_tree) => p_tree,
            Err(e) => panic!("failed to revparse: {}", e),
        };
        let mut diff_opts = DiffOptions::new();
        diff_opts.ignore_whitespace(true);
        diff_opts.context_lines(0);
        diff_opts.ignore_filemode(true);
        // diff_opts.indent_heuristic(true);
        let diff = repo.diff_tree_to_tree(Some(&p_tree), Some(&commit_tree), Some(&mut diff_opts)).unwrap();

        // println!("Diff: {:#?}", diff.unwrap().print(DiffFormat::Patch, print_callback).unwrap());
        let diff_stats = diff.stats().unwrap();
        println!("Diff stats: {} insertions, {} deletions", diff_stats.insertions(), diff_stats.deletions() );

        println!("--- DIFF ---");
        diff.foreach(&mut file_cb, Some(&mut binary_cb), Some(&mut hunk_cb), Some(&mut line_cb));
    }
}

fn file_cb(d:DiffDelta, progress:f32 ) -> bool {
    true
}
fn binary_cb(d:DiffDelta, db:DiffBinary ) -> bool {
    true
}
fn hunk_cb(d:DiffDelta, dh:DiffHunk) -> bool {
    true
}
fn line_cb(d:DiffDelta, o:Option<DiffHunk>, l:DiffLine ) -> bool {
    print!("{}   {}", l.origin(), String::from_utf8(l.content().to_vec()).unwrap());
    true
}

#[allow(unused_variables)]
fn print_callback(d:DiffDelta, o:Option<DiffHunk>, l:DiffLine ) -> bool {
    true
}

#[cfg(test)]
mod tests {

    #[test]
    fn test_goofing_tag() {
        let repo = match git2::Repository::open(".") {
            Ok(repo) => repo,
            Err(e) => panic!("failed to open: {}", e),
        };
        let head = match repo.revparse_single("test-goofing") {
            Ok(head) => head,
            Err(e) => panic!("failed to revparse: {}", e),
        };

        assert_eq!("79caa008ba1f9d06b34b4acc7c03d7fade185a63", head.id().to_string());
    }
}
