import { parseResponse } from 'scripts/parse-response';
import { writable, get, derived } from 'svelte/store';
import { GetDateFormats, GetSettings, SaveSettingsGui } from 'wailsjs/go/main/App'

interface Settings {
	Version: string;
  Workflow: string;
  Theme: string;
  DateFormat: number;
  HighlightMainBranch: boolean;
  HighlightConventionalCommits: boolean;
	DisplayCommitSignatureInList: boolean;
  DisplayAvatars: boolean;
  BackgroundOpacity: number;
}

function createSettings() {
  const { subscribe, set, update } = writable({} as Settings);

  return {
    subscribe,
    refresh: async () => {
      await GetSettings().then(result => {
        set(result);
      });
    },
    updateWorkflow: (workflow: string) => {
      if (!["squash", "rebase"].includes(workflow)) {
        workflow = "merge";
      }
      update(s => {
        s.Workflow = workflow;
        return s;
      });
      settings.save();
    },
    updateTheme: (theme: string) => {
      update(s => {
        s.Theme = theme;
        return s;
      });
      settings.save();
    },
    updateDateFormat: (format: number) => {
      update(s => {
        s.DateFormat = format;
        return s;
      });
      settings.save();
    },
    toggleHighlightConventionalCommits: () => {
      update(s => {
        s.HighlightConventionalCommits = !s.HighlightConventionalCommits;
        return s;
      });
      settings.save();
    },
    toggleHighlightMainBranch: () => {
      update(s => {
        s.HighlightMainBranch = !s.HighlightMainBranch;
        return s;
      });
      settings.save();
    },
    toggleDisplayAvatars: () => {
      update(s => {
        s.DisplayAvatars = !s.DisplayAvatars;
        return s;
      });
      settings.save();
    },
    save: async () => {
      SaveSettingsGui(get(settings)).then(result => {
        parseResponse(result);
      });
    }
  };
}
export const settings = createSettings();

export const theme = derived(settings, $settings => {
  if (!$settings.Theme) {
    // If no theme is set, see if user prefers light mode.
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
      // Do not set Theme, it should remain unset.
      return "light";
    }
    // Default to dark.
    return "dark";
  }
  return $settings.Theme;
});

type DateFormats = {[id: number]: {
  Display: string;
  // format: string;
}}

function createDateFormats() {
  const { subscribe, set } = writable({} as DateFormats);

  return {
    subscribe,
    fetch: async () => {
      if (!Object.keys(get(dateFormats)).length) {
        await GetDateFormats().then(result => {
          set(result);
        });
      }
      return get(dateFormats);
    },
  };
}
export const dateFormats = createDateFormats();
dateFormats.fetch(); // only once for the app
