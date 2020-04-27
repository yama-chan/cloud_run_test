export GOOGLE_APPLICATION_CREDENTIALS={サービスアカウントキーのパス}
export GOOGLE_CLOUD_PROJECT={プロジェクト名}

cd cloud_run_test
docker build ./ -t test_cloud_run
docker run -p 8080:8080 test_cloud_run
http://localhost:8080/test
