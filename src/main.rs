use git2::Repository;

fn main() {
    let repo = match Repository::open(".") {
        Ok(repo) => repo,
        Err(e) => panic!("failed to open: {}", e),
    };
    let head = match repo.revparse_single("HEAD") {
        Ok(head) => head,
        Err(e) => panic!("failed to revparse: {}", e),
    };
    println!("Repo head is: {}",head.id());
}
