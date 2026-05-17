// Exemplo de command Tauri. Para invocar do Angular:
//   import { invoke } from '@tauri-apps/api/core'
//   await invoke<string>('greet', { name: 'mundo' })
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Olá, {}!", name)
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![greet])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
