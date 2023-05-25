import { EventsOn } from 'wailsjs/runtime/runtime';

export function registerConflictEvents() {
  EventsOn('rebase-conflicts', (e) => {
    console.log('todo: rebase conflicts');
  });

  EventsOn('merge-conflicts', (e) => {
    console.log('todo: merge conflicts');
  });

  EventsOn('revert-conflicts', (e) => {
    console.log('todo: revert conflicts');
  });

  EventsOn('unresolved-conflicts', (e) => {
    console.log('todo: unresolved conflicts');
  });
}
