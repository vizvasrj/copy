import time
import os
import shutil
import multiprocessing

def copy_item(item):
    source_item, destination_item = item
    try:
        if os.path.isfile(source_item):
            shutil.copy2(source_item, destination_item)
        else:
            shutil.copytree(source_item, destination_item)
        # print(f"Copied item: {destination_item}")
    except Exception as e:
        print(f"Failed to copy item: {e}")

if __name__ == '__main__':
    start_time = time.time()
    # Configuration
    source_path = "demo_files4/"
    destination_path = "/home/tmp"


    # Create a multiprocessing.Pool
    with multiprocessing.Pool() as pool:
        # Iterate over files and directories in source path
        for root, dirs, files in os.walk(source_path):
            for file in files:
                source_file = os.path.join(root, file)
                destination_file = os.path.join(
                    destination_path, 
                    os.path.relpath(source_file, source_path)
                )
                item = (source_file, destination_file)
                pool.apply_async(copy_item, (item,))
            for dir in dirs:
                source_dir = os.path.join(root, dir)
                destination_dir = os.path.join(
                    destination_path, 
                    os.path.relpath(source_dir, source_path)
                )
                item = (source_dir, destination_dir)
                pool.apply_async(copy_item, (item,))

        # Close the pool and wait for all processes to finish
        pool.close()
        pool.join()

    print("All items copied and took seconds =", 
          time.time() - start_time)
