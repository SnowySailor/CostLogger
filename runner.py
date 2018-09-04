import time
import subprocess
import sys
import os
import signal
from watchdog.observers import Observer
from watchdog.events import PatternMatchingEventHandler

class MyHandler(PatternMatchingEventHandler):
    patterns           = ["*.go", "*.template", "*.yaml"]
    current_process    = None

    def process(self, event):
        print(event.src_path + " was " + event.event_type)
        self.rebuild_and_run()

    def on_modified(self, event):
        self.process(event)

    def on_created(self, event):
        self.process(event)

    def on_moved(self, event):
        self.process(event)

    def on_deleted(self, event):
        self.process(event)

    def rebuild_and_run(self):
        self.kill_current_process()
        print("Rebuilding...")
        success = self.run_build()
        if success:
            print("Running...")
            self.run_process()

    def kill_current_process(self):
        if self.current_process is None:
            return
        try:
            os.kill(self.current_process.pid, signal.SIGTERM)
        except Exception:
            print("Unable to stop current process. Restart this tool.")

    def run_build(self):
        try:
            build_process = subprocess.Popen(['/bin/sh', 'build.sh'], stderr=subprocess.PIPE)
            _, stderr = build_process.communicate()
            error = stderr.decode('utf-8')
            if error != "":
                print(error)
                return False
            return True
        except Exception as e:
            print("Unable to rebuild project: " + str(e))

    def run_process(self):
        try:
            self.current_process = subprocess.Popen(['./site'], cwd='site')
        except Exception as e:
            print("Unable to start process: " + str(e))

def main():
    args = sys.argv[1:]
    observer = Observer()
    handler = MyHandler()
    handler.rebuild_and_run()
    observer.schedule(handler, path=args[0] if args else '.', recursive=True)
    observer.start()

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()

    observer.join()

if __name__ == '__main__':
    main()