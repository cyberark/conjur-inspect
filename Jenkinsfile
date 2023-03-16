#!/usr/bin/env groovy

// Automated release, promotion and dependencies
properties([
  // Include the automated release parameters for the build
  release.addParams(),
])

// Performs release promotion.  No other stages will be run
if (params.MODE == "PROMOTE") {
  release.promote(params.VERSION_TO_PROMOTE) { sourceVersion, targetVersion, assetDirectory ->
    // Any assets from sourceVersion Github release are available in assetDirectory
    // Any version number updates from sourceVersion to targetVersion occur here
    // Any publishing of targetVersion artifacts occur here
    // Anything added to assetDirectory will be attached to the Github Release
  }
  return
}

pipeline {
  agent { label 'executor-v2' }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
    timeout(time: 2, unit: 'HOURS')
  }

  triggers {
    cron(getDailyCronString())
  }

  environment {
    // Sets the MODE to the specified or autocalculated value as appropriate
    MODE = release.canonicalizeMode()
  }

  stages {
    // Aborts any builds triggered by another project that wouldn't include any changes
    stage ("Skip build if triggering job didn't create a release") {
      when {
        expression {
          MODE == "SKIP"
        }
      }
      steps {
        script {
          currentBuild.result = 'ABORTED'
          error("Aborting build because this build was triggered from upstream, but no release was built")
        }
      }
    }

    // Generates a VERSION file based on the current build number and latest version in CHANGELOG.md
    stage('Validate Changelog and set version') {
      steps {
        sh './bin/parse-changelog'
        updateVersion("CHANGELOG.md", "${BUILD_NUMBER}")
      }
    }

    stage('Get latest upstream dependencies') {
      steps {
        updateGoDependencies('${WORKSPACE}/go.mod')
      }
    }


    stage('Run Unit Tests') {
      steps {
        sh './bin/test-unit'
      }

      post {
        always {
          junit 'junit.xml'

          cobertura autoUpdateHealth: false,
            autoUpdateStability: false,
            coberturaReportFile: 'coverage.xml',
            conditionalCoverageTargets: '70, 0, 0',
            failUnhealthy: false,
            failUnstable: false,
            maxNumberOfBuilds: 0,
            lineCoverageTargets: '70, 0, 0',
            methodCoverageTargets: '70, 0, 0',
            onlyStable: false,
            sourceEncoding: 'ASCII',
            zoomCoverageChart: false
            ccCoverage("gocov", "--prefix github.com/cyberark/conjur-inspect")
        }
      }
    }

    // This produces the conjur-inspect binaries for integration tests and
    // pushing a release when this is a RELEASE build.
    stage('Create Release Assets') {
      steps {
        sh "bin/build-release"
      }
    }

    // Currently the integration tests don't pass or fail the build based on
    // any conditions. Rather, that provide an easy way to inspect the result
    // from a few various environments without running the tool manually.
    stage('Run Integration Tests') {
      steps {
        sh 'bin/test-integration'
      }
      post {
        always {
           archiveArtifacts artifacts: 'ci/integration/results/**', allowEmptyArchive: true, fingerprint: false
        }
      }
    }

    stage('Release') {
      when {
        expression {
          MODE == "RELEASE"
        }
      }

      steps {
        release { billOfMaterialsDirectory, assetDirectory, toolsDirectory ->
          // Publish release artifacts to all the appropriate locations
          // Copy any artifacts to assetDirectory to attach them to the Github release

          // Create Go application SBOM using the go.mod version for the golang container image
          sh """go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --main "cmd/conjur-inspect/" --output "${billOfMaterialsDirectory}/go-app-bom.json" """
          // Create Go module SBOM
          sh """go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --output "${billOfMaterialsDirectory}/go-mod-bom.json" """

          // Add goreleaser artifacts to release
          sh """cp dist/*.tar.gz "${assetDirectory}" """
          sh """cp dist/*.rpm "${assetDirectory}" """
          sh """cp dist/*.deb "${assetDirectory}" """
          sh """cp "dist/SHA256SUMS.txt" "${assetDirectory}" """
        }
      }
    }
  }
}
