use std::io::Write;
use std::{io, path::PathBuf};

use aingle_types::prelude::{
    AppBundle, AppManifest, AppManifestCurrentBuilder, AppSlotManifest, SafBundle, SafManifest,
};

fn readline(prompt: Option<&str>) -> io::Result<Option<String>> {
    let mut input = String::new();
    if let Some(prompt) = prompt {
        print!("{} ", prompt);
        io::stdout().flush()?;
    }
    io::stdin().read_line(&mut input)?;
    let input = input.trim();
    Ok(if input.is_empty() {
        None
    } else {
        Some(input.to_owned())
    })
}

fn prompt_default<S: Into<String>>(prompt: &str, default: S) -> io::Result<String> {
    let default = default.into();
    let prompt = format!("{} ({})", prompt, default);
    Ok(readline(Some(&prompt))?.unwrap_or(default))
}

fn prompt_optional(prompt: &str) -> io::Result<Option<String>> {
    readline(Some(prompt))
}

fn prompt_required(prompt: &str) -> io::Result<String> {
    loop {
        if let Some(line) = readline(Some(prompt))? {
            return Ok(line);
        }
    }
}

fn prompt_saf_init(root_dir: PathBuf) -> anyhow::Result<SafBundle> {
    let name = prompt_required("name:")?;
    let uid = Some(prompt_default(
        "uid:",
        "00000000-0000-0000-0000-000000000000",
    )?);
    let manifest = SafManifest::current(name, uid, None, vec![]);
    Ok(SafBundle::new(manifest, vec![], root_dir)?)
}

fn prompt_app_init(root_dir: PathBuf) -> anyhow::Result<AppBundle> {
    let name = prompt_required("name:")?;
    let description = prompt_optional("description:")?;
    let slot = AppSlotManifest::sample("sample-slot".into());
    let manifest: AppManifest = AppManifestCurrentBuilder::default()
        .name(name)
        .description(description)
        .slots(vec![slot])
        .build()
        .unwrap()
        .into();
    Ok(mr_bundle::Bundle::new(manifest, vec![], root_dir)?.into())
}

pub async fn init_saf(target: PathBuf) -> anyhow::Result<()> {
    let bundle = prompt_saf_init(target.to_owned())?;
    bundle.unpack_yaml(&target, false).await?;
    Ok(())
}

pub async fn init_app(target: PathBuf) -> anyhow::Result<()> {
    let bundle = prompt_app_init(target.to_owned())?;
    bundle.unpack_yaml(&target, false).await?;
    Ok(())
}

#[cfg(test)]
mod tests {

    // TODO: make these functions able to take an arbitrary stream so that
    //       they can be tested

    // use super::*;

    // #[tokio::test]
    // async fn can_init_saf() {
    //     let tmpdir = tempdir::TempDir::new("ai_bundle").unwrap();
    //     init_saf(tmpdir.path().join("app")).await.unwrap();
    //     init_saf(tmpdir.path().join("app/n")).await.unwrap();
    // }
}
