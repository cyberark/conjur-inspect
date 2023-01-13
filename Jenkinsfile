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
        sh './bin/parse_changelog'
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
            ccCoverage("gocov", "--prefix github.com/conjurinc/conjur-preflight")
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
        // Build release artifacts
        sh "bin/build_release"

        release { billOfMaterialsDirectory, assetDirectory, toolsDirectory ->
          // Publish release artifacts to all the appropriate locations
          // Copy any artifacts to assetDirectory to attach them to the Github release

          // Create Go application SBOM using the go.mod version for the golang container image
          sh """go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --main "cmd/conjur-preflight/" --output "${billOfMaterialsDirectory}/go-app-bom.json" """
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
