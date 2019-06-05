use actix_web::{web, App, HttpResponse, HttpServer, Responder};
use serde_derive::Deserialize;
use git2::*;

mod git_churn;

#[derive(Deserialize)]
struct ChurnArgs {
    repo: String,
    commit: String,
}

fn churn(args: web::Query<ChurnArgs>) -> HttpResponse {
    let repo = Repository::open(&args.repo).expect("Could not open repo");
    let commit = repo
        .revparse_single(&args.commit)
        .expect("Failed revparse")
        .peel_to_commit()
        .expect("Could not peel to commit");
    let stats = git_churn::compute_churn(&repo, &commit, false);
    return HttpResponse::Ok().json(stats);
}

fn main() -> std::io::Result<()> {
    HttpServer::new(|| App::new().service(web::resource("/churn").to(churn)))
        .bind("127.0.0.1:8080")?
        .run()
}
