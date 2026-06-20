import * as vscode from 'vscode';
import { execFile } from 'child_process';

// Map a VS Code languageId to a per-line comment prefix for bigmark's -c flag
// (used by tier 1/2/3 and rotated modes). Default is "# ".
const PREFIX_BY_LANG: Record<string, string> = {
  javascript: '// ',
  javascriptreact: '// ',
  typescript: '// ',
  typescriptreact: '// ',
  c: '// ',
  cpp: '// ',
  csharp: '// ',
  go: '// ',
  rust: '// ',
  java: '// ',
  kotlin: '// ',
  swift: '// ',
  scala: '// ',
  php: '// ',
  css: '/* ', // single-line /* is unusual but keeps it valid-ish; canvas mode is the real CSS path
  scss: '// ',
  less: '// ',
  python: '# ',
  ruby: '# ',
  perl: '# ',
  shellscript: '# ',
  yaml: '# ',
  toml: '# ',
  dockerfile: '# ',
  makefile: '# ',
  r: '# ',
  sql: '-- ',
  lua: '-- ',
  haskell: '-- ',
  html: '<!-- ',
  xml: '<!-- ',
};

// Map a VS Code languageId to one of bigmark's supported --lang tokens for
// canvas mode (canvasDelims in internal/render/render.go). Anything not listed
// falls back to "js" (which bigmark also defaults to).
const CANVAS_LANG_BY_LANG: Record<string, string> = {
  javascript: 'js',
  javascriptreact: 'js',
  typescript: 'ts',
  typescriptreact: 'ts',
  php: 'php',
  css: 'css',
  scss: 'css',
  less: 'css',
  html: 'html',
  xml: 'html',
  python: 'python',
  ruby: 'ruby',
};

function prefixForLanguage(languageId: string): string {
  return PREFIX_BY_LANG[languageId] ?? '# ';
}

function canvasLangForLanguage(languageId: string): string {
  return CANVAS_LANG_BY_LANG[languageId] ?? 'js';
}

function buildArgs(
  mode: string,
  width: number,
  languageId: string,
  configuredPrefix: string,
  label: string
): string[] {
  const args: string[] = [];

  switch (mode) {
    case 'tier2':
      args.push('-2');
      break;
    case 'tier3':
      args.push('-3');
      break;
    case 'rotated':
      args.push('-r');
      break;
    case 'canvas':
      args.push('--canvas');
      break;
    case 'tier1':
    default:
      break;
  }

  args.push('-w', String(width));

  if (mode === 'canvas') {
    args.push('--lang', canvasLangForLanguage(languageId));
  } else {
    const prefix = configuredPrefix !== '' ? configuredPrefix : prefixForLanguage(languageId);
    args.push('-c', prefix);
  }

  args.push(label);
  return args;
}

function runBigmark(
  binaryPath: string,
  args: string[],
  cwd: string | undefined,
  figletPath: string
): Promise<string> {
  const env = { ...process.env };
  if (figletPath !== '') {
    env.BIGMARK_FIGLET = figletPath;
  }
  return new Promise((resolve, reject) => {
    execFile(binaryPath, args, { cwd, env }, (err, stdout, stderr) => {
      if (err) {
        const detail = (stderr && stderr.trim()) || err.message;
        reject(new Error(detail));
        return;
      }
      resolve(stdout);
    });
  });
}

export function activate(context: vscode.ExtensionContext) {
  const disposable = vscode.commands.registerCommand('bigmark.insert', async () => {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
      vscode.window.showErrorMessage('bigmark: open a file first.');
      return;
    }

    const selection = editor.selection;
    const selectedText = editor.document.getText(selection).trim();

    const label = await vscode.window.showInputBox({
      prompt: 'bigmark banner label',
      value: selectedText,
      placeHolder: 'SECTION',
    });
    if (label === undefined || label.trim() === '') {
      return; // cancelled or empty
    }

    const cfg = vscode.workspace.getConfiguration('bigmark');
    const binaryPath = cfg.get<string>('binaryPath', 'bigmark');
    const mode = cfg.get<string>('mode', 'tier1');
    const width = cfg.get<number>('width', 80);
    const commentPrefix = cfg.get<string>('commentPrefix', '');
    const figletPath = cfg.get<string>('figletPath', '');

    const languageId = editor.document.languageId;
    const args = buildArgs(mode, width, languageId, commentPrefix, label);

    const cwd = vscode.workspace.getWorkspaceFolder(editor.document.uri)?.uri.fsPath;

    let output: string;
    try {
      output = await runBigmark(binaryPath, args, cwd, figletPath);
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e);
      vscode.window.showErrorMessage(`bigmark: ${msg}`);
      return;
    }

    // strip a single trailing newline so the banner sits on its own lines
    const banner = output.replace(/\n$/, '');

    await editor.edit((b) => {
      if (!selection.isEmpty) {
        b.replace(selection, banner);
      } else {
        b.insert(selection.active, banner);
      }
    });
  });

  context.subscriptions.push(disposable);
}

export function deactivate() {}
