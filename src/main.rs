use git2::*;
use std::*;
use futures::executor::*;

mod git_churn;

async fn async_main() {
    let repo = Repository::open(".").expect("Could not open repo");
    let commit = repo
        .revparse_single("test-first-test")
        .expect("Failed revparse")
        .peel_to_commit()
        .expect("Could not peel to commit");
    let stats = git_churn::compute_churn(&repo, &commit, false);
    println!("{:#?}", stats);
}

fn main() {
    let future = async_main();
    block_on(future);
    println!("Done!!");
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
